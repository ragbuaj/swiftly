export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  full_name: string;
  username: string;
  phone_number: string;
}

export interface AuthResponse {
  status: string;
  message: string;
  data: {
    access_token: string;
    refresh_token: string;
  };
}

export interface GoogleLoginRequest {
  id_token: string;
}

export interface ForgotPasswordRequest {
  identifier: string; // email, username, or phone
}

export interface ForgotPasswordResponse {
  status: string;
  message: string;
  data?: {
    token: string;
    email_hint: string;
  };
}

export interface ResetPasswordRequest {
  token: string;
  new_password: string;
}
