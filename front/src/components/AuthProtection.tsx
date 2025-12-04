import React from "react";
import useAuthStore from "../stores/authStore";
import Login from "./Login";

type Props = {
  children?: React.ReactNode;
};

const AuthProtection: React.FC<Props> = ({ children }) => {
  const isAuthenticate = useAuthStore((state) => state.isAuthenticated);

  // if isAuthenticated() is false, redirect to login page
  if (!isAuthenticate) {
    return <Login />;
  }

  return children;
};

export default AuthProtection;
