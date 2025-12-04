import { create } from "zustand";

type Me = {
  id: string;
  name?: string;
  email?: string;
  avatarUrl?: string;
  [key: string]: any;
};

type AuthState = {
  me: Me | null;
  setMe: (me: Me | null) => void;
  clearMe: () => void;
  isAuthenticated: () => boolean;
};

const useAuthStore = create<AuthState>((set, get) => ({
  me: null,
  setMe: (me: Me | null) => set({ me }),
  clearMe: () => set({ me: null }),
  isAuthenticated: () => {
    const me = get().me;
    return me !== null && !!me.id;
  },
}));

export default useAuthStore;
