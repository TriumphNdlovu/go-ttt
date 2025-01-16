const board = document.querySelector(".board");
const cells = document.querySelectorAll(".cell");
const turnIndicator = document.querySelector(".player-turn");
const player1Info = document.querySelector(".player-1");
const player2Info = document.querySelector(".player-2");

let currentPlayer = 'X'; // Player 1 starts
let gameBoard = ['', '', '', '', '', '', '', '', '']; // Game state tracking

// WebSocket connection for multiplayer
const ws = new WebSocket("ws://localhost:8080"); // Update with your server's URL

// Handle incoming messages
ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    if (data.board) {
        gameBoard = data.board;
        renderBoard();
    }
    if (data.turn) {
        currentPlayer = data.turn;
        updateTurnIndicator();
};

// Render the board state
function renderBoard() {
    gameBoard.forEach((cell, index) => {
        cells[index].textContent = cell;
        if (cell === 'X') {
            cells[index].style.color = "#1e90ff"; // Blue for Player 1
        } else if (cell === 'O') {
            cells[index].style.color = "#ff6347"; // Red for Player 2
        } else {
            cells[index].style.color = "#000000"; // Default color for empty
        }
    });
}

// Handle cell click (player move)
cells.forEach(cell => {
    cell.addEventListener("click", () => {
        const index = cell.dataset.index;
        if (!gameBoard[index]) {
            gameBoard[index] = currentPlayer;
            ws.send(JSON.stringify({ board: gameBoard, turn: currentPlayer === 'X' ? 'O' : 'X' }));
            renderBoard();
        }
    });
});

// Update the turn indicator
function updateTurnIndicator() {
    if (currentPlayer === 'X') {
        turnIndicator.textContent = "Player 1's turn";
        turnIndicator.classList.add('player-1');
        turnIndicator.classList.remove('player-2');
    } else {
        turnIndicator.textContent = "Player 2's turn";
        turnIndicator.classList.add('player-2');
        turnIndicator.classList.remove('player-1');
    }
}

// Initial rendering
renderBoard();
updateTurnIndicator();
}