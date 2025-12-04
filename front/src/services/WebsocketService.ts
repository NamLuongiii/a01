// // Kết nối WebSocket
// const ws = new WebSocket("ws://localhost:8080/ws?user_id=123&room_id=room1");

// // Lắng nghe messages
// ws.onmessage = function (event) {
//   const data = JSON.parse(event.data);
//   console.log("Received:", data);
// };

// // Gửi message
// const sendMessage = () => {
//   ws.send(
//     JSON.stringify({
//       type: "chat_message",
//       content: "Hello World!",
//       room_id: "room1",
//     })
//   );
// };

// // Join room khác
// const joinRoom = (roomId: number) => {
//   ws.send(
//     JSON.stringify({
//       type: "join_room",
//       room_id: roomId,
//     })
//   );
// };
