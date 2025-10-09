package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	db "github.com/mike-woub/User_Auth/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

// func (apiCfg *apiConfig) handlerSignup(w http.ResponseWriter, r *http.Request) {
// 	type signupParams struct {
// 		Username string `json:"username"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}
// 	var params signupParams
// 	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
// 		http.Error(w, "invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		http.Error(w, "failed to hash password", http.StatusInternalServerError)
// 		return
// 	}

// 	user, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
// 		Username: params.Username,
// 		Email:    params.Email,
// 		Password: string(hashedPassword),
// 	})
// 	if err != nil {
// 		log.Printf("CreateUser error: %v", err)
// 		http.Error(w, "failed to create user", http.StatusInternalServerError)
// 		return
// 	}

// 	response := struct {
// 		Username string `json:"username"`
// 		Email    string `json:"email"`
// 	}{
// 		Username: user.Username,
// 		Email:    user.Email,
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(response)
// }

func (apiCfg *apiConfig) handlerSignup(w http.ResponseWriter, r *http.Request) {
	type signupParams struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	var params signupParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	params.Email = strings.ToLower(params.Email)

	// Check if email or username already exists
	if count, _ := apiCfg.DB.CheckEmailExists(r.Context(), db.CheckEmailExistsParams{Email: params.Email, ID: 0}); count > 0 {
		http.Error(w, "email already in use", http.StatusConflict)
		return
	}
	if count, _ := apiCfg.DB.CheckUsernameExists(r.Context(), db.CheckUsernameExistsParams{Username: params.Username, ID: 0}); count > 0 {
		http.Error(w, "username already taken", http.StatusConflict)
		return
	}

	// Generate OTP
	otp := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	apiCfg.OTPStore.Set(params.Email, otp, 5*time.Minute)

	// Send OTP via Mailtrap
	body := fmt.Sprintf("Hello %s,\n\nYour verification code is: %s\n\nThanks,\nLocalLoop Team", params.Username, otp)
	if err := sendEmail(params.Email, "Your OTP Code", body); err != nil {
		http.Error(w, "failed to send verification email", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Verification code sent to email"})
}

func (apiCfg *apiConfig) handlerVerifySignup(w http.ResponseWriter, r *http.Request) {
	type verifyParams struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		OTP      string `json:"otp"`
	}
	var params verifyParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	params.Email = strings.ToLower(params.Email)

	storedOTP, found := apiCfg.OTPStore.Get(params.Email)
	if !found || storedOTP != params.OTP {
		http.Error(w, "invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		Username: params.Username,
		Email:    params.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Signup successful",
		"username": user.Username,
		"email":    user.Email,
	})
}
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type loginResponse struct {
		Token string `json:"token"`
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Normalize email to lowercase
	req.Email = strings.ToLower(req.Email)

	user, err := apiCfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar_url": func() string {
			if user.AvatarUrl.Valid {
				return user.AvatarUrl.String
			}
			return ""
		}(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := loginResponse{Token: signedToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (apiCfg *apiConfig) handlerUpdateUsername(w http.ResponseWriter, r *http.Request) {
	type request struct {
		NewUsername string `json:"new_username"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := apiCfg.DB.CheckUsernameExists(r.Context(), db.CheckUsernameExistsParams{
		Username: req.NewUsername,
		ID:       userID,
	})
	if err != nil {
		http.Error(w, "failed to check username", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "username already taken", http.StatusConflict)
		return
	}

	err = apiCfg.DB.UpdateUsername(r.Context(), db.UpdateUsernameParams{
		Username: req.NewUsername,
		ID:       userID,
	})
	if err != nil {
		http.Error(w, "failed to update username", http.StatusInternalServerError)
		return
	}

	user, err := apiCfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch updated user", http.StatusInternalServerError)
		return
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Username updated successfully",
		"token":   signedToken,
	})
}

func (apiCfg *apiConfig) handlerUpdateEmail(w http.ResponseWriter, r *http.Request) {
	type request struct {
		NewEmail string `json:"new_email"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := apiCfg.DB.CheckEmailExists(r.Context(), db.CheckEmailExistsParams{
		Email: req.NewEmail,
		ID:    userID,
	})
	if err != nil {
		http.Error(w, "failed to check email", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "email already in use", http.StatusConflict)
		return
	}

	err = apiCfg.DB.UpdateEmail(r.Context(), db.UpdateEmailParams{
		Email: req.NewEmail,
		ID:    userID,
	})
	if err != nil {
		http.Error(w, "failed to update email", http.StatusInternalServerError)
		return
	}

	user, err := apiCfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch updated user", http.StatusInternalServerError)
		return
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email updated successfully",
		"token":   signedToken,
	})
}

func (apiCfg *apiConfig) handlerUpdatePassword(w http.ResponseWriter, r *http.Request) {
	type request struct {
		NewPassword string `json:"new_password"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	err = apiCfg.DB.UpdatePassword(r.Context(), db.UpdatePasswordParams{
		Password: string(hashedPassword),
		ID:       userID,
	})
	if err != nil {
		http.Error(w, "failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

func (apiCfg *apiConfig) handlerUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	type request struct {
		AvatarURL string `json:"avatar_url"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = apiCfg.DB.UpdateAvatar(r.Context(), db.UpdateAvatarParams{
		AvatarUrl: sql.NullString{String: req.AvatarURL, Valid: req.AvatarURL != ""},
		ID:        userID,
	})
	if err != nil {
		http.Error(w, "failed to update avatar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Avatar updated successfully"})
}

func (apiCfg *apiConfig) handlerRequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Email string `json:"email"`
	}
	var params requestParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	params.Email = strings.ToLower(params.Email)

	// Check if user exists
	_, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Generate OTP
	otp := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	apiCfg.OTPStore.Set(params.Email, otp, 5*time.Minute)

	// Send OTP via email
	if err := sendEmail(params.Email, "Password Reset Code", otp); err != nil {
		http.Error(w, "failed to send OTP", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent to email"})
}
func (apiCfg *apiConfig) handlerVerifyPasswordReset(w http.ResponseWriter, r *http.Request) {
	type verifyParams struct {
		Email    string `json:"email"`
		OTP      string `json:"otp"`
		Password string `json:"new_password"`
	}
	var params verifyParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	params.Email = strings.ToLower(params.Email)

	storedOTP, found := apiCfg.OTPStore.Get(params.Email)
	if !found || storedOTP != params.OTP {
		http.Error(w, "invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	err = apiCfg.DB.UpdateUserPassword(r.Context(), db.UpdateUserPasswordParams{
		Email:    params.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, "failed to update password", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successful"})
}
