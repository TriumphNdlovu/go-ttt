let myPlayer = null; // The role assigned to this player ("X" or "O")
let currentPlayer = null; // The current player's turn ("X" or "O")
let gameBoard = Array(9).fill(""); // Represents the Tic-Tac-Toe board
let gameOver = false; // Tracks if the game is over

// Select the turn indicator element
const turnIndicator = document.getElementById("turn-indicator");

// WebSocket setup
const ws = new WebSocket("ws://localhost:8080");

ws.onopen = () => {
    console.log("Connected to the server.");
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);

    if (data.error) {
        alert(data.error);
        return;
    }

    // Assign the player's role if not already assigned
    if (data.player && myPlayer === null) {
        myPlayer = data.player;
        alert(`You are player ${myPlayer}`);
    }

    // Update the game board
    if (data.board) {
        gameBoard = data.board;
        renderBoard();
    }

    // Update the current turn
    if (data.turn) {
        currentPlayer = data.turn;
        updateTurnIndicator();
    }
};

// Render the board
function renderBoard() {
    const cells = document.querySelectorAll(".cell");
    cells.forEach((cell, index) => {
        cell.textContent = gameBoard[index];
    });
}

// Update the turn indicator
function updateTurnIndicator() {
    if (gameOver) {
        turnIndicator.textContent = "Game Over!";
        return;
    }
    if (currentPlayer === myPlayer) {
        turnIndicator.textContent = "Your turn!";
    } else {
        turnIndicator.textContent = `Waiting for ${currentPlayer}'s turn...`;
    }
}

// Add click event listeners to the cells
function initializeBoard() {
    const cells = document.querySelectorAll(".cell");
    cells.forEach((cell, index) => {
        cell.addEventListener("click", () => {
            // Check if it's the player's turn and if the game is ongoing
            if (gameOver || currentPlayer !== myPlayer || gameBoard[index]) return;

            // Make the move
            gameBoard[index] = myPlayer;

            // Send the move to the server
            ws.send(JSON.stringify({
                board: gameBoard,
                player: myPlayer,
            }));

            renderBoard();
        });
    });
}

// Call the function to initialize the board
initializeBoard();
