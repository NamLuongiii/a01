import { create } from "zustand";
import type { User } from "../types";

interface Me extends User {}

type AuthState = {
  me: Me | null;
  isAuthenticated: boolean;

  setMe: (me: Me | null) => void;
  clearMe: () => void;
};

const useAuthStore = create<AuthState>((set, get) => ({
  me: null,
  isAuthenticated: false,
  setMe: (me: Me | null) => set({ me, isAuthenticated: !!me }),
  clearMe: () => set({ me: null, isAuthenticated: false }),
}));

export default useAuthStore;
