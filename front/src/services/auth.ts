import type { ApiResponse, User } from "../types";
import api from "./api";

export const AuthService = {
  login: async (name: string) => {
    // Implement login logic here
    const res = await api.post<
      ApiResponse<User> & {
        user?: User;
      }
    >("/auth/login", { name });
    return res.data;
  },
};
