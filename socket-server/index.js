import { Server } from "socket.io";
import Redis from "ioredis";

const io = new Server({ cors: "*" });
const subscriber = new Redis({
  username: "default",
  password: "SdfNsULIdJ3ZQivu4I9Uvr2owZCwFPhO",
  port: 17256,
  host: "redis-17256.c330.asia-south1-1.gce.redns.redis-cloud.com",
});
const PORT = 9002;

io.on("connection", (socket) => {
  socket.on("subscribe", (channel) => {
    socket.join(channel);
    socket.emit("message", `Joined ${channel}`);
  });
});

io.listen(PORT, () => console.log(`Socket Server ${PORT}`));

async function initRedisSubscribe() {
  console.log("Subscribed to logs....");
  subscriber.psubscribe("logs:*");
  subscriber.on("pmessage", (pattern, channel, message) => {
    io.to(channel).emit("message", message);
  });
}

initRedisSubscribe();
