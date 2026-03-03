export interface User {
  id: string;
  email: string;
  username?: string;
  phone_number?: string;
  full_name: string;
  role: string;
  status: string;
  email_verified_at?: string;
  deleted_at?: string;
  date_of_birth?: string;
  gender?: string;
  newsletter_subscribed: boolean;
  avatar_url?: string;
  bio?: string;
  default_shipping_address_id?: string;
  default_billing_address_id?: string;
  created_at: string;
  updated_at: string;
}

export interface UpdateProfileRequest {
  full_name: string;
  username: string;
  phone_number: string;
  bio: string;
  gender: string;
  date_of_birth?: string;
  newsletter_subscribed: boolean;
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
