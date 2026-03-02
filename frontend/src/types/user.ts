export interface User {
  id: string;
  email: string;
  username?: string;
  phone_number?: string;
  full_name: string;
  created_at: string;
  updated_at: string;
}

export interface UserProfileResponse {
  status: string;
  message: string;
  data: User;
}

export interface CheckEmailResponse {
  status: string;
  message: string;
  data: {
    available: boolean;
  };
}
