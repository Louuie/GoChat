package page

import (
	"net/http"
)

func main() {
	// Serve HTML frontend
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<meta name="htmx-config" content='{"selfRequestsOnly":false}'>
				<title>GoChat</title>
				<script src="https://unpkg.com/htmx.org@2.0.4"></script>
				<script>
					document.addEventListener("htmx:afterSwap", (event) => {
					    // Check if the swap target was the #rooms-container
					    if (event.detail.target.id === "rooms-container") {
					        // Attach click event listener to each button
					        const buttons = document.querySelectorAll(".room-btn");
					        buttons.forEach((button) => {
					            button.addEventListener("click", () => {
					                // Extract the UUID from the clicked button
					                const uuid = button.getAttribute("data-uuid");
					                console.log("Room UUID clicked:", uuid);

					                // Open a WebSocket connection for the clicked UUID
					                const ws = new WebSocket('ws://localhost:8080/ws?uuid=${uuid}');

					                ws.onopen = () => {
					                    console.log('WebSocket connection opened for UUID: ${uuid}');
					                };

					                ws.onmessage = (event) => {
					                    console.log('Message from server: ${event.data}');
					                };

					                ws.onclose = () => {
					                    console.log('WebSocket connection closed for UUID: ${uuid}');
					                };
					            });
					        });
					    }
					});
					</script>

			</head>
			<body>
				<h1>GoChat</h1>
				<form onsubmit="handleUsername(event)" id="username-form">
					<input type="text" id="username-input" placeholder="Type your username" required>
					<button type="submit">Send</button>
				</form>
				<div id="whole-chat" style="visibility: hidden;">
					<div id="chat" style="border: 1px solid #ccc; padding: 10px; height: 300px; overflow-y: scroll;"></div>
					<form onsubmit="sendMessage(event)">
						<input type="text" id="messageInput" placeholder="Type your message" required>
						<button type="submit">Send</button>
					</form>
				</div>
				<div id="rooms-container">
					<button
	    				hx-get="http://localhost:5050/fetch-rooms"
					    hx-target="#rooms-container"
					    hx-swap="innerHTML">
					    Fetch All Rooms
					</button>
				</div>
			</body>
			</html>
		`))
	})

	// Start HTTP server on port 3030
	http.ListenAndServe(":3030", nil)
}

func StartPage() {
	main()
}
