### Notes

Server -> Client
1. Welcomes a new client. Asks for a name. If no valid name given, disconnects the client. Informative  
error messages sent to the client are allowed.

Server -> All Clients
1. When receiving a message from one client, publishes messages to all clients.
2. When a client enters and sets their name, informs all clients that the person has entered the room.
3. When a client leaves, informs all clients that the person has left the room.


Client -> Server
1. Informs the server its name. 
2. Sends messages to the server