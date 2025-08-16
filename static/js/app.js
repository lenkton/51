function roll(gameId, dice) {
    event.preventDefault();
    fetch("/games/"+gameId+"/roll",
        { method:"POST", body: JSON.stringify({ dice: parseInt(dice, 10) }) })
        .then((response) => response.json())
        .then((v) => console.log(v));
}
function joinGame(gameId, userName) {
    event.preventDefault();
    fetch("/games/"+gameId+"/join",
        { method:"POST", body: JSON.stringify({ userName: userName }) })
        .then((response) => response.json())
        .then((v) => console.log(v));
}
function createGame() {
    fetch("/games", { method: "POST" })
        .then((response) => {
            if (response.ok) {
                response.json().then((message)=>{
                    document.location = "/games/" + message.id;
                });
            };
        });
}
function startGame(gameId) {
    fetch(`/games/${gameId}/start`, { method: "POST" })
        .then((response) => {
            if (!response.ok) {
                response.json().then((message) => {
                    console.log(message.error);
                });
            };
        });
}

var players = [];
function addPlayer(player) {
    players.push(player);
    turns[player.id] = [];

    const playerList = document.getElementById('players-list');

    let li = document.createElement("li");
    li.appendChild(document.createTextNode(player.name));
    playerList.appendChild(li);

    renderTurnsTable()
}
// should be set in the HTML
var turns = {};
// TODO: do not rerender the whole table every time an update occurs
function renderTurnsTable() {
    let newTable = document.createElement('table');
    let thead = document.createElement('thead');
    let tbody = document.createElement('tbody');
    newTable.appendChild(thead);
    newTable.appendChild(tbody);

    let roundsCount = 0;
    for (const player of players) {
        const playerTurns = turns[player.id].length;
        if (playerTurns > roundsCount) {
            roundsCount = playerTurns;
        }
    }
    for (let i = 0; i < roundsCount; i++) {
        let row = tbody.insertRow();
        for (const player of players) {
            const turn = turns[player.id][i];
            const newCell = row.insertCell();
            if (turn != undefined) {
                newCell.appendChild(document.createTextNode(turn.result))
            }
        }
    }
    const namesRow = thead.insertRow();
    for (const player of players) {
        let playerTotal = 0;
        if (turns[player.id]) {
            for (const turn of turns[player.id]) {
                playerTotal += turn.result;
            }
        }
        const nameCell = namesRow.appendChild(document.createElement('th'));
        nameCell.appendChild(document.createTextNode(`${player.name} (${playerTotal})`))
    }

    const gameTable = document.getElementById('turns');
    newTable.id = 'turns';
    gameTable.replaceWith(newTable);
}
function addTurn(turn, player) {
    turns[player.id].push(turn);
    renderTurnsTable()
}

const NameTotalRE = /(.*) \((\d*)\)$/;
function updatePlayerTotal(playerIndex, toAdd) {
    const gameTable = document.getElementById('turns');
    let nameCell = gameTable.rows[0].cells[playerIndex];
    let nameAndTotal = nameCell.innerText;
    const [name, total] = NameTotalRE.exec(nameAndTotal).slice(1);
    const newTotal = +total + +toAdd;
    nameCell.innerText = `${name} (${newTotal})`;
}

function updateHeading(game) {
    let heading = document.getElementById('game-heading');
    heading.innerText = `Game ${game.id} (${game.status})`;
}
function connectToGameUpdates(gameId) {
    const url = "/games/" + gameId + "/updates";
    ws = new WebSocket(url);
    ws.onmessage = (event) => {
        console.log("I have received an update from the server: "+event.data);
        if (event.data[0] == '{') {
            const msg = JSON.parse(event.data);
            switch (msg.type) {
                case "newPlayer":
                    addPlayer(msg.player);
                    break;
                case "newTurn":
                    addTurn(msg.turn, msg.player);
                    break;
                case "gameStarted":
                    updateHeading(msg.game);
                    break;
            }
        }
    };
    ws.onopen = (event) => {
        ws.send("12312");
    };
}
