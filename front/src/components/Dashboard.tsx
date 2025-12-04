import React from "react";
import useAuthStore from "../stores/authStore";

const Dashboard: React.FC = () => {
  const me = useAuthStore((state) => state.me);

  return (
    <div>
      <header className="flex items-center gap-2 justify-end p-4">
        <div className="size-8 rounded-full bg-green-200"></div>
        <div>Welcome, {me?.name}</div>
      </header>
    </div>
  );
};

export default Dashboard;
