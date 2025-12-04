import api from "./api";

export const RoomService = {
  getAll: () => api.get("/rooms"),
  getById: (id: string) => api.get(`/rooms/${id}`),
  create: (data: any) => api.post("/rooms", data),
  update: (id: string, data: any) => api.put(`/rooms/${id}`, data),
  delete: (id: string) => api.delete(`/rooms/${id}`),
};
