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
function addPlayer(player) {
    const playerList = document.getElementById('players-list');
    const gameTable = document.getElementById('turns');

    let li = document.createElement("li");
    li.appendChild(document.createTextNode(player.name));
    playerList.appendChild(li);

    const rows = gameTable.getElementsByTagName("tr");
    // WARN: i assume, we have only one row (only the heading)
    for (let row of rows) {
        // todo: if i is 0 then th, otherwise - td
        let th = document.createElement("th");
        th.appendChild(document.createTextNode(player.name + " (0)"));
        row.appendChild(th);
    };
}
function addTurn(turn) {
    const gameTable = document.getElementById('turns');

    const playersCount = gameTable.rows[0].cells.length;

    const rows = gameTable.rows;
    const lastTurnInRound = rows[rows.length-1].cells.length == playersCount;
    if (lastTurnInRound) {
        gameTable.insertRow(-1);
    }
    const lastRow = rows[rows.length-1];
    let td = lastRow.insertCell(-1);
    td.appendChild(document.createTextNode(turn.result));

    updatePlayerTotal(lastRow.cells.length-1, turn.result);
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
                    addTurn(msg.turn);
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
