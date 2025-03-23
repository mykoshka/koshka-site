# Koshka API Documentation

## Overview
Koshka API provides a backend service for managing users, pets, and administrative tasks. This documentation details every available endpoint, including expected inputs, outputs, and possible error messages.

---

## Authentication Endpoints

### **Login**
`POST /login`

#### **Description:**
Authenticates a user and returns a JWT token.

#### **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### **Responses:**
- **200 OK**
```json
{
  "token": "jwt-token-here"
}
```
- **400 Bad Request**
```json
{
  "error": "Invalid email or password"
}
```

### **Register**
`POST /register`

#### **Description:**
Registers a new user.

#### **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}
```

#### **Responses:**
- **201 Created**
```json
{
  "message": "User registered successfully"
}
```
- **400 Bad Request**
```json
{
  "error": "Email already in use"
}
```

### **Request Password Reset**
`POST /request-password-reset`

#### **Description:**
Initiates a password reset process by sending a reset link to the user’s email.

#### **Request Body:**
```json
{
  "email": "user@example.com"
}
```

#### **Responses:**
- **200 OK**
```json
{
  "message": "Password reset link sent"
}
```
- **404 Not Found**
```json
{
  "error": "Email not found"
}
```

### **Reset Password**
`POST /reset-password`

#### **Description:**
Completes the password reset process using a reset token.

#### **Request Body:**
```json
{
  "token": "reset-token-here",
  "new_password": "newsecurepassword"
}
```

#### **Responses:**
- **200 OK**
```json
{
  "message": "Password reset successfully"
}
```
- **400 Bad Request**
```json
{
  "error": "Invalid or expired token"
}
```

### **Verify Email**
`GET /verify-email`

#### **Description:**
Verifies a user’s email address using a verification token.

#### **Query Parameters:**
- `token`: Verification token sent to the user’s email (required)

#### **Responses:**
- **200 OK**
```json
{
  "message": "Email verified successfully"
}
```
- **400 Bad Request**
```json
{
  "error": "Invalid or expired token"
}
```

---

## User Profile Endpoints

### **Get User Profile**
`GET /api/profile`

#### **Description:**
Fetches the profile details of the authenticated user.

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Responses:**
- **200 OK**
```json
{
  "email": "user@example.com",
  "name": "John Doe",
  "address": "123 Main St",
  "mobile_number": "555-1234",
  "joined_on": "2024-01-01",
  "pets": [
    {
      "id": "123",
      "name": "Buddy",
      "neutered": true,
      "vaccinated": false,
      "date_of_birth": "2020-05-10"
    }
  ]
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

### **Change Password**
`POST /api/profile/auth/change-password`

#### **Description:**
Allows the authenticated user to change their password.

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Request Body:**
```json
{
  "old_password": "securepassword",
  "new_password": "newsecurepassword"
}
```

#### **Responses:**
- **200 OK**
```json
{
  "message": "Password changed successfully"
}
```
- **400 Bad Request**
```json
{
  "error": "Incorrect old password"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

---

## Pet Management Endpoints

### **Add Pet**
`POST /api/pet/add`

#### **Description:**
Adds a new pet to the authenticated user’s profile.

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Request Body:**
```json
{
  "name": "Buddy",
  "neutered": true,
  "vaccinated": false,
  "date_of_birth": "2020-05-10"
}
```

#### **Responses:**
- **201 Created**
```json
{
  "message": "Pet added successfully",
  "pet_id": "123"
}
```
- **400 Bad Request**
```json
{
  "error": "Invalid pet details"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

### **Update Pet**
`PUT /api/pet/:id`

#### **Description:**
Updates details of an existing pet.

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Request Body:**
```json
{
  "name": "Buddy",
  "neutered": true,
  "vaccinated": true,
  "date_of_birth": "2020-05-10"
}
```

#### **Responses:**
- **200 OK**
```json
{
  "message": "Pet updated successfully"
}
```
- **404 Not Found**
```json
{
  "error": "Pet not found"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

### **Delete Pet**
`DELETE /api/pet/:id`

#### **Description:**
Deletes a pet from the user’s profile.

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Responses:**
- **200 OK**
```json
{
  "message": "Pet deleted successfully"
}
```
- **404 Not Found**
```json
{
  "error": "Pet not found"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

### **Assign Collar to Pet**
`POST /api/pets/:id/assign-collar`

#### **Description:**
Assigns a collar to a specific pet.

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Request Body:**
```json
{
  "tag_id": "TAG123"
}
```

#### **Responses:**
- **200 OK**
```json
{
  "message": "Collar assigned successfully"
}
```
- **404 Not Found**
```json
{
  "error": "Pet or collar not found"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

---

## Admin Endpoints

### **List All Users**
`GET /api/admin/ListUsers`

#### **Description:**
Retrieves a list of all registered users (admin only).

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Responses:**
- **200 OK**
```json
{
  "users": [
    {
      "email": "user@example.com",
      "name": "John Doe",
      "joined_on": "2024-01-01"
    }
  ]
}
```
- **403 Forbidden**
```json
{
  "error": "Admin access required"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

### **List API Keys**
`GET /api/admin/api_keys`

#### **Description:**
Retrieves a list of all API keys (admin only).

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Responses:**
- **200 OK**
```json
{
  "api_keys": [
    {
      "key": "abc123",
      "created_at": "2024-01-01",
      "permissions": ["read", "write"]
    }
  ]
}
```
- **403 Forbidden**
```json
{
  "error": "Admin access required"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

### **Create API Key**
`POST /api/admin/api_keys`

#### **Description:**
Creates a new API key (admin only).

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Request Body:**
```json
{
  "permissions": ["read", "write"]
}
```

#### **Responses:**
- **201 Created**
```json
{
  "message": "API key created",
  "key": "abc123"
}
```
- **403 Forbidden**
```json
{
  "error": "Admin access required"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

### **Update API Key Permissions**
`POST /api/admin/api_keys/permissions/:key`

#### **Description:**
Updates permissions for a specific API key (admin only).

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Request Body:**
```json
{
  "permissions": ["read", "write", "delete"]
}
```

#### **Responses:**
- **200 OK**
```json
{
  "message": "Permissions updated successfully"
}
```
- **404 Not Found**
```json
{
  "error": "API key not found"
}
```
- **403 Forbidden**
```json
{
  "error": "Admin access required"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

### **List API Routes**
`GET /api/admin/routes`

#### **Description:**
Retrieves a list of all API routes available (admin only).

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Responses:**
- **200 OK**
```json
{
  "routes": ["/api/profile", "/api/pet/:id", "/api/admin/api_keys"]
}
```
- **403 Forbidden**
```json
{
  "error": "Admin access required"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

### **Upload Tags**
`POST /api/admin/upload-tags`

#### **Description:**
Uploads multiple pet tags in bulk (admin only).

#### **Headers:**
```
Authorization: Bearer <token>
```

#### **Request Body:**
```json
{
  "tags": [
    {
      "tag_id": "TAG123",
      "status": "Active"
    },
    {
      "tag_id": "TAG124",
      "status": "Inactive"
    }
  ]
}
```

#### **Responses:**
- **200 OK**
```json
{
  "message": "Tags uploaded successfully"
}
```
- **400 Bad Request**
```json
{
  "error": "Invalid tag data"
}
```
- **403 Forbidden**
```json
{
  "error": "Admin access required"
}
```
- **401 Unauthorized**
```json
{
  "error": "Missing or invalid token"
}
```

### **Fetch Phone Number by Code**
`GET /api/v1/reunite/ext/:code`

#### **Description:**
Fetches the phone number associated with a lost pet’s tag code.

#### **Responses:**
- **200 OK**
```json
{
  "phone_number": "+1234567890"
}
```
- **404 Not Found**
```json
{
  "error": "Tag code not found"
}
```

### **Fetch Collar Information**
`GET /collar/:tag_id`

#### **Description:**
Retrieves information about a specific pet collar.

#### **Responses:**
- **200 OK**
```json
{
  "collar_id": "TAG123",
  "assigned_pet": "Buddy",
  "status": "Active"
}
```
- **404 Not Found**
```json
{
  "error": "Collar not found"
}
```

### **Static Image Handling**
`HEAD /images`

#### **Description:**
Handles static image requests.

#### **Responses:**
- **200 OK** (Image exists)
- **404 Not Found** (Image does not exist)

---

## Security
This API requires authentication using JWT tokens for most endpoints. Always include the following header in authenticated requests:
```
Authorization: Bearer <your_token>
```

---

## Notes
- Ensure all requests are properly formatted as JSON.
- Some endpoints require admin privileges.
- Tokens expire after 24 hours and must be refreshed.

---

🚀 **Koshka API - Secure & Scalable Pet Management Backend**
