import { useQuery } from "@tanstack/react-query";
import React, { useEffect } from "react";
import { QUERY_KEYS } from "../constants";
import { RoomService } from "../services";
import useAuthStore from "../stores/authStore";

const Dashboard: React.FC = () => {
  const me = useAuthStore((state) => state.me);
  const clear = useAuthStore((state) => state.clearMe);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");
    ws.onopen = () => console.log("Connected!");
    ws.onmessage = (event) => console.log("Received:", event.data);
    ws.send("Hello from client!");
  }, []);

  // get rooms and display them
  const { data: rooms, isLoading } = useQuery({
    queryKey: [QUERY_KEYS.ROOMS],
    queryFn: async () => {
      const rooms = await RoomService.getAll();
      return rooms;
    },
    throwOnError: true,
  });

  return (
    <div>
      <header className="flex items-center gap-2 justify-end p-4">
        <div
          className="size-8 rounded-full bg-green-200 cursor-pointer"
          onClick={clear}
        ></div>
        <div>Welcome, {me?.name}</div>
      </header>

      <div>
        {/* Display Loading  */}
        {isLoading && <div>Loading rooms...</div>}

        {/* Display Rooms */}
        {rooms &&
          rooms.map((room) => (
            <div key={room.id} className="p-4 border-b">
              <h2 className="text-lg font-bold">{room.name}</h2>
              <p>{room.description}</p>
            </div>
          ))}
      </div>
    </div>
  );
};

export default Dashboard;
