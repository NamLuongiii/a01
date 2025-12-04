export interface User {
  id: number;
  name: string;
  email: string;
  // other user fields
}

export interface ApiResponse<T> {
  data: T;
  message: string;
  success: boolean;
}

export interface Room {
  id: number;
  name: string;
  users: User[];
  description: string;
}
