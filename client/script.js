const socket = new WebSocket('ws://localhost:8080');

let playerId = null;
let selectedCharacter = null;
let selectedCharacterRow = null;
let selectedCharacterCol = null;
let placementArray = [];

const gameBoard = document.getElementById('game-board');
const currentPlayerSpan = document.getElementById('current-player');
const moveHistoryList = document.getElementById('move-list');
const gameOverScreen = document.getElementById('game-over-screen');
const gameOverMessage = document.getElementById('game-over-message');
const newGameButton = document.getElementById('new-game-button');
const moveButtonsDiv = document.getElementById('move-buttons');

socket.onopen = () => {
    const playerName = prompt("Enter your name:");
    sendMessage('join_game', { playerName });
};

socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    switch (message.type) {
        case 'game_setup':
            playerId = message.data;
            placementArray = [];
            alert(`You are Player ${playerId}`);
            break;
        case 'game_start':
            updateGameBoard(message.data);
            updateCurrentPlayer(message.data.CurrentTurn);
            break;
        case 'game_state':
            updateGameBoard(message.data);
            updateCurrentPlayer(message.data.CurrentTurn);
            updateMoveHistory(message.data);
            break;
        case 'error':
            alert(message.data);
            break;
        case 'game_over':
            showGameOverScreen(message.data);
            break;
    }
};

socket.onclose = () => {
    console.log('Websocket connection closed');
};

function sendMessage(type, data) {
    socket.send(JSON.stringify({ type, data }));
}

function createGameBoard() {
    for (let row = 0; row < 5; row++) {
        for (let col = 0; col < 5; col++) {
            const cell = document.createElement('div');
            cell.classList.add('cell');
            cell.id = `cell-${row}-${col}`;
            cell.addEventListener('click', handleCellClick);
            gameBoard.appendChild(cell);
        }
    }
}

function handleCellClick(event) {
    const cellId = event.target.id;
    const [_, row, col] = cellId.split('-').map(Number);

    if (playerId === null) {
        return;
    }

    if (placementArray.length < 5) {
        const startingRow = playerId === 1 ? 4 : 0; 
        if (row === startingRow) {
            const characterType = prompt("Enter character type (P1, P2, P3, P4, P5, H1, or H2):");
            if (characterType) {
                placementArray.push(characterType);
                event.target.textContent = `${playerId}-${characterType}`;
                event.target.classList.add(getPlayerClass(playerId));

                if (placementArray.length === 5) {
                    sendMessage('setup_done', placementArray);
                }
            }
        }
    } else {
        const cellContent = event.target.textContent;
        if (cellContent && cellContent.startsWith(`${playerId}-`)) {
            selectedCharacter = cellContent.substring(2);
            selectedCharacterRow = row;
            selectedCharacterCol = col;
            highlightValidMoves(row, col, selectedCharacter);
        } else if (selectedCharacter && event.target.classList.contains('highlighted')) {
            const moveDirection = getMoveDirection(row, col);
            sendMessage('make_move', `${selectedCharacter}:${moveDirection}`);
            selectedCharacter = null;
            selectedCharacterRow = null;
            selectedCharacterCol = null;
            clearHighlights(); 
        }
    }
}

function highlightValidMoves(row, col, character) {
    clearHighlights();

    switch (character) {
        case 'P1':
        case 'P2':
        case 'P3':
        case 'P4':
        case 'P5':
            highlightPawnMoves(row, col);
            break;
        case 'H1':
            highlightHero1Moves(row, col);
            break;
        case 'H2':
            highlightHero2Moves(row, col);
            break;
    }
}

function highlightPawnMoves(row, col) {
    highlightCellIfValid(row - 1, col);
    highlightCellIfValid(row + 1, col);
    highlightCellIfValid(row, col - 1);
    highlightCellIfValid(row, col + 1);
}

function highlightHero1Moves(row, col) {
    highlightCellIfValid(row - 2, col);
    highlightCellIfValid(row + 2, col);
    highlightCellIfValid(row, col - 2);
    highlightCellIfValid(row, col + 2);
}

function highlightHero2Moves(row, col) {
    highlightCellIfValid(row - 2, col - 2);
    highlightCellIfValid(row - 2, col + 2);
    highlightCellIfValid(row + 2, col - 2);
    highlightCellIfValid(row + 2, col + 2);
}

function highlightCellIfValid(row, col) {
    if (row >= 0 && row < 5 && col >= 0 && col < 5) {
        const cell = document.getElementById(`cell-${row}-${col}`);
        cell.classList.add('highlighted');
    }
}

function clearHighlights() {
    const highlightedCells = document.querySelectorAll('.highlighted');
    highlightedCells.forEach(cell => cell.classList.remove('highlighted'));
}

function getMoveDirection(targetRow, targetCol) {
    const rowDiff = targetRow - selectedCharacterRow;
    const colDiff = targetCol - selectedCharacterCol;

    if (rowDiff === -1 && colDiff === 0) return 'F';
    if (rowDiff === 1 && colDiff === 0) return 'B';
    if (rowDiff === 0 && colDiff === -1) return 'L';
    if (rowDiff === 0 && colDiff === 1) return 'R';
    if (rowDiff === -2 && colDiff === 0) return 'F';
    if (rowDiff === 2 && colDiff === 0) return 'B';
    if (rowDiff === 0 && colDiff === -2) return 'L';
    if (rowDiff === 0 && colDiff === 2) return 'R';
    if (rowDiff === -2 && colDiff === -2) return 'FL';
    if (rowDiff === -2 && colDiff === 2) return 'FR';
    if (rowDiff === 2 && colDiff === -2) return 'BL';
    if (rowDiff === 2 && colDiff === 2) return 'BR';

    return '';
}

function updateGameBoard(gameState) {
    gameBoard.innerHTML = ''; 
    createGameBoard(); 

    for (let row = 0; row < 5; row++) {
        for (let col = 0; col < 5; col++) {
            const charCode = gameState.Board[row][col];
            if (charCode) {
                const cell = document.getElementById(`cell-${row}-${col}`);
                const charElement = document.createElement('div');
                charElement.textContent = charCode;
                charElement.classList.add(getPlayerClass(charCode[0]));
                cell.appendChild(charElement);
            }
        }
    }
}

function updateCurrentPlayer(turn) {
    currentPlayerSpan.textContent = `Current Player: Player ${turn + 1}`;
}

function updateMoveHistory(gameState) {
    moveHistoryList.innerHTML = ''; 

    if (gameState.MoveHistory && gameState.MoveHistory.length > 0) {
        gameState.MoveHistory.forEach(move => {
            const moveItem = document.createElement('li');
            moveItem.textContent = move;
            moveHistoryList.appendChild(moveItem);
        });
    }
}

function showGameOverScreen(message) {
    gameOverMessage.textContent = message;
    gameOverScreen.classList.remove('hidden');
}

newGameButton.addEventListener('click', () => {
    sendMessage('new_game');
    gameOverScreen.classList.add('hidden'); 
    gameBoard.innerHTML = ''; 
    createGameBoard(); 
    moveHistoryList.innerHTML = ''; 
    playerId = null; 
    selectedCharacter = null; 
    selectedCharacterRow = null;
    selectedCharacterCol = null;
    placementArray = []; 
});

function getPlayerClass(playerId) {
    return playerId === 'A' ? 'player-a' : 'player-b';
}

createGameBoard();