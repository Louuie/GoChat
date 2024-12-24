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
    // Establish a WebSocket connection
    const ws = new WebSocket("ws://localhost:8080/ws?uuid=some-uuid");

    let lastSentMessage = ""; // Store the last message sent by the client

    // Handle username submission
    function handleUsername(event) {
        event.preventDefault(); // Prevent the form from refreshing the page

        // Get the username input and chat elements
        const usernameInput = document.getElementById("username-input");
        const wholeChat = document.getElementById("whole-chat");
        const usernameForm = document.getElementById("username-form");

        // Check if the username is valid
        if (usernameInput.value.trim() !== "") {
            wholeChat.style.visibility = "visible"; // Show the chat interface
            usernameForm.style.visibility = "hidden"; // Hide the username form
        }
    }

    // Handle sending a message
    function sendMessage(event) {
        event.preventDefault(); // Prevent form submission default behavior
        const input = document.getElementById("messageInput");
        const message = input.value.trim();

        if (message !== "") {
            ws.send(message); // Send the message to the server
            lastSentMessage = message; // Store the message
            const chatDiv = document.getElementById("chat");
            const usernameInput = document.getElementById("username-input").value;
            chatDiv.insertAdjacentHTML(
                "beforeend",
                '<div>${usernameInput}: ${message}</div>'
            ); // Display the user's message
            input.value = ""; // Clear the input field
        }
    }

    // Handle incoming messages
    ws.onmessage = (event) => {
        const chatDiv = document.getElementById("chat");

        // Avoid displaying duplicate messages
        if (event.data !== lastSentMessage) {
            chatDiv.insertAdjacentHTML(
                "beforeend",
                '<div>Other: ${event.data}</div>'
            );
        }
    };

    // Handle dynamic content loaded via HTMX
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

                    // Define sendMessage function specific to this WebSocket
                    function sendMessage(event) {
                        event.preventDefault(); // Prevent form submission default behavior
                        const input = document.getElementById("messageInput");
                        const message = input.value.trim();

                        if (message !== "") {
                            ws.send(message); // Send the message to the server
                            const chatDiv = document.getElementById("chat");
                            const usernameInput = document.getElementById("username-input").value;
                            chatDiv.insertAdjacentHTML(
                                "beforeend",
                                '<div>${usernameInput}: ${message}</div>'
                            ); // Display the user's message
                            input.value = ""; // Clear the input field
                        }
                    }

                    // Handle incoming messages from the WebSocket
                    ws.onmessage = (event) => {
                        const chatDiv = document.getElementById("chat");

                        // Avoid displaying duplicate messages
                        if (event.data !== lastSentMessage) {
                            chatDiv.insertAdjacentHTML(
                                "beforeend",
                                '<div>Other: ${event.data}</div>'
                            );
                        }
                    };

                    // Attach the sendMessage logic to the send button
                    document
                        .getElementById("sendButton")
                        .addEventListener("click", sendMessage);
                });
            });
        }
    });
</script>



			</head>
			<body>
				<h1>GoChat</h1>
				<div id="whole-chat">
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
