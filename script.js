// Replace with your actual Railway backend URL
const BACKEND_URL = "https://poketactix-production.up.railway.app";

// Handle search form on index.html
document.addEventListener("DOMContentLoaded", function() {
    const searchForm = document.getElementById("searchForm");
    if (searchForm) {
        searchForm.addEventListener("submit", function(e) {
            e.preventDefault();
            const input = document.getElementById("searchInput");
            const name = input.value.trim().toLowerCase();
            if (name) {
                window.location.href = `pokemon.html?name=${encodeURIComponent(name)}`;
            }
        });
    }

    // Load pokemon.html data
    if (window.location.pathname.endsWith("pokemon.html")) {
        const params = new URLSearchParams(window.location.search);
        const name = params.get("name");
        if (name) {
            fetch(`${BACKEND_URL}/pokemon?name=${encodeURIComponent(name)}`)
                .then(res => res.json())
                .then(data => {
                    if (data.error) {
                        document.getElementById("pokemonCard").innerHTML = `<p>${data.error}</p>`;
                    } else {
                        // Render card (customize as needed)
                        document.getElementById("pokemonCard").innerHTML = `
                            <div class="pokemon-header">
                                <h2>${data.Name}</h2>
                                <div class="types">${data.Types.map(type => `<span class="type ${type}">${type}</span>`).join('')}</div>
                            </div>
                            <div class="pokemon-image">
                                <img src="${data.Sprite}" alt="${data.Name}">
                            </div>
                            <div class="stats">
                                <div class="stat"><span class="stat-label">HP</span><div class="stat-bar"><div class="stat-fill" style="width: ${(data.HP/255)*100}%"></div></div><span class="stat-value">${data.HP}</span></div>
                                <div class="stat"><span class="stat-label">Stamina</span><div class="stat-bar"><div class="stat-fill" style="width: ${(data.Stamina/255)*100}%"></div></div><span class="stat-value">${data.Stamina}</span></div>
                                <div class="stat"><span class="stat-label">Attack</span><div class="stat-bar"><div class="stat-fill" style="width: ${(data.Attack/255)*100}%"></div></div><span class="stat-value">${data.Attack}</span></div>
                                <div class="stat"><span class="stat-label">Defense</span><div class="stat-bar"><div class="stat-fill" style="width: ${(data.Defense/255)*100}%"></div></div><span class="stat-value">${data.Defense}</span></div>
                                <div class="stat"><span class="stat-label">Speed</span><div class="stat-bar"><div class="stat-fill" style="width: ${(data.Speed/255)*100}%"></div></div><span class="stat-value">${data.Speed}</span></div>
                            </div>
                            <div class="moves">
                                <h3>Moves</h3>
                                <div class="moves-list">${data.Moves.map(move => `<span class="move">${move.name} (Power: ${move.power}, Type: ${move.attack_type})</span>`).join('')}</div>
                            </div>
                            <a href="index.html" class="btn back-btn">Back to Home</a>
                        `;
                    }
                })
                .catch(err => {
                    document.getElementById("pokemonCard").innerHTML = `<p>Error: ${err}</p>`;
                });
        }
    }

    // Load battle-arena.html data
    if (window.location.pathname.endsWith("battle-arena.html")) {
        let session = null;
        let state = null;
        let battleOver = false;
        let currentTurn = 1;
        let isPlayerTurn = true;
        let playerName = "Pikachu";
        let aiName = "Charizard";
        const params = new URLSearchParams(window.location.search);
        if (params.get("player")) playerName = params.get("player");
        if (params.get("ai")) aiName = params.get("ai");

        // Player elements
        const playerNameEl = document.getElementById("playerName");
        const playerTypesEl = document.getElementById("playerTypes");
        const playerSpriteEl = document.getElementById("playerSprite");
        const playerHP = document.getElementById("playerHP");
        const playerHPMax = document.getElementById("playerHPMax");
        const playerHPBar = document.getElementById("playerHPBar");
        const playerStamina = document.getElementById("playerStamina");
        const playerStaminaMax = document.getElementById("playerStaminaMax");
        const playerStaminaBar = document.getElementById("playerStaminaBar");
        const playerAttack = document.getElementById("playerAttack");
        const playerDefense = document.getElementById("playerDefense");
        const playerSpeed = document.getElementById("playerSpeed");

        // AI elements
        const aiNameEl = document.getElementById("aiName");
        const aiTypesEl = document.getElementById("aiTypes");
        const aiSpriteEl = document.getElementById("aiSprite");
        const aiHP = document.getElementById("aiHP");
        const aiHPMax = document.getElementById("aiHPMax");
        const aiHPBar = document.getElementById("aiHPBar");
        const aiStamina = document.getElementById("aiStamina");
        const aiStaminaMax = document.getElementById("aiStaminaMax");
        const aiStaminaBar = document.getElementById("aiStaminaBar");
        const aiAttack = document.getElementById("aiAttack");
        const aiDefense = document.getElementById("aiDefense");
        const aiSpeed = document.getElementById("aiSpeed");

        // Other elements
        const battleLogDiv = document.getElementById("battleLog");
        const battleModeStatus = document.getElementById("battleModeStatus");
        const attackBtn = document.getElementById("attackBtn");
        const defendBtn = document.getElementById("defendBtn");
        const passBtn = document.getElementById("passBtn");
        const sacrificeBtn = document.getElementById("sacrificeBtn");
        const surrenderBtn = document.getElementById("surrenderBtn");
        const movesModal = document.getElementById("movesModal");
        const movesGrid = document.getElementById("movesGrid");
        const cancelMoveBtn = document.getElementById("cancelMoveBtn");
        const battleCommands = document.getElementById("battleCommands");

        let logHistory = [];
        let whoseTurn = "player";

        function renderTypes(types) {
            if (!types) return "";
            return types.map(type => `<span class="pokemon-type type-${type.toLowerCase()}">${type}</span>`).join(' ');
        }

        function updateCard(card, isPlayer) {
            if (!card) return;
            if (isPlayer) {
                playerNameEl.textContent = card.Name;
                playerTypesEl.innerHTML = renderTypes(card.Types);
                playerSpriteEl.src = card.Sprite;
                playerSpriteEl.alt = card.Name;
                playerHP.textContent = card.HP;
                playerHPMax.textContent = card.HPMax || card.HP;
                playerHPBar.style.width = ((card.HP / (card.HPMax || card.HP)) * 100) + "%";
                playerStamina.textContent = card.Stamina;
                playerStaminaMax.textContent = (card.Speed || 1) * 2;
                playerStaminaBar.style.width = ((card.Stamina / ((card.Speed || 1) * 2)) * 100) + "%";
                playerAttack.textContent = card.Attack;
                playerDefense.textContent = card.Defense;
                playerSpeed.textContent = card.Speed;
            } else {
                aiNameEl.textContent = card.Name;
                aiTypesEl.innerHTML = renderTypes(card.Types);
                aiSpriteEl.src = card.Sprite;
                aiSpriteEl.alt = card.Name;
                aiHP.textContent = card.HP;
                aiHPMax.textContent = card.HPMax || card.HP;
                aiHPBar.style.width = ((card.HP / (card.HPMax || card.HP)) * 100) + "%";
                aiStamina.textContent = card.Stamina;
                aiStaminaMax.textContent = (card.Speed || 1) * 2;
                aiStaminaBar.style.width = ((card.Stamina / ((card.Speed || 1) * 2)) * 100) + "%";
                aiAttack.textContent = card.Attack;
                aiDefense.textContent = card.Defense;
                aiSpeed.textContent = card.Speed;
            }
        }

        function renderLog(logArr) {
            if (!logArr || !logArr.length) return "";
            return logArr.map(entry => `<div class="log-entry">${entry}</div>`).join("");
        }

        function updateUI() {
            if (!state) return;
            updateCard(state.Player.Deck[state.PlayerActiveIdx], true);
            updateCard(state.AI.Deck[state.AIActiveIdx], false);
            // Log
            battleLogDiv.innerHTML = renderLog(logHistory);
            // Status
            battleModeStatus.textContent = `Turn ${state.TurnNumber || 1} - ${whoseTurn === "player" ? "Player's" : "AI's"} Turn`;
            // Enable/disable buttons based on whose turn
            [attackBtn, defendBtn, passBtn, sacrificeBtn, surrenderBtn].forEach(btn => {
                btn.disabled = (whoseTurn !== "player") || battleOver;
            });
            if (battleOver) {
                let resultMsg = "Battle Over";
                if (state.winner === "player") {
                    resultMsg = "You won!";
                } else if (state.winner === "ai") {
                    resultMsg = "You lost";
                } else if (state.winner === "draw") {
                    resultMsg = "Draw!";
                }
                battleCommands.innerHTML = `
                    <div style="text-align: center; color: white; font-size: 1.5rem; font-weight: bold; width: 100%;">
                        🎉 ${resultMsg} 🎉
                        <br><br>
                        <button class="command-btn attack-btn" onclick="location.reload()" style="margin-right: 1rem;">Play Again</button>
                        <button class="command-btn defend-btn" onclick="window.location.href='battle.html'">New Battle</button>
                    </div>
                `;
            }
        }

        function appendLog(entries) {
            if (!Array.isArray(entries)) entries = [entries];
            logHistory.push(...entries);
            battleLogDiv.innerHTML = renderLog(logHistory);
            battleLogDiv.scrollTop = battleLogDiv.scrollHeight;
        }

        function startBattle() {
            fetch(`${BACKEND_URL}/battle/start`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ playerName, aiName })
            })
            .then(res => res.json())
            .then(data => {
                session = data.session;
                state = data.state;
                state.winner = data.winner;
                whoseTurn = data.whoseTurn || "player";
                logHistory = ["The battle begins! Choose your action."];
                battleOver = false;
                updateUI();
                if (whoseTurn === "ai") {
                    handleAITurn();
                }
            });
        }

        function handleAITurn() {
            // Disable player controls
            [attackBtn, defendBtn, passBtn, sacrificeBtn, surrenderBtn].forEach(btn => btn.disabled = true);
            // POST to backend to let AI act (no move needed from player)
            fetch(`${BACKEND_URL}/battle/move`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ session, move: "" })
            })
            .then(res => res.json())
            .then(data => {
                state = data.state;
                state.winner = data.winner;
                whoseTurn = data.whoseTurn || "player";
                if (data.log) appendLog(data.log);
                if (data.battleOver || state.BattleOver) {
                    battleOver = true;
                }
                updateUI();
                // If after AI move it's still AI's turn (shouldn't happen), repeat
                if (!battleOver && whoseTurn === "ai") {
                    setTimeout(handleAITurn, 800);
                }
            });
        }

        function doMove(move, moveIdx) {
            if (!session || battleOver || whoseTurn !== "player") return;
            let payload = { session, move };
            if (move === "attack") payload.moveIdx = moveIdx;
            fetch(`${BACKEND_URL}/battle/move`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload)
            })
            .then(res => res.json())
            .then(data => {
                state = data.state;
                state.winner = data.winner;
                whoseTurn = data.whoseTurn || "player";
                if (data.log) appendLog(data.log);
                if (data.battleOver || state.BattleOver) {
                    battleOver = true;
                }
                updateUI();
                if (!battleOver && whoseTurn === "ai") {
                    setTimeout(handleAITurn, 800);
                }
            });
        }

        attackBtn.onclick = () => {
            // Show moves modal
            const moves = state.Player.Deck[state.PlayerActiveIdx].Moves;
            movesGrid.innerHTML = moves.map((move, i) => {
                // Normalize type name
                const type = (move.Type || move.attack_type || '').toLowerCase().replace(/\s+/g, '');
                const moveName = move.Name || move.name || 'Move';
                const movePower = move.Power || move.power || '?';
                return `
                    <button class="move-btn" data-move-idx="${i}">
                        <div class="move-name">${moveName}</div>
                        <span class="pokemon-type type-${type}">${type.charAt(0).toUpperCase() + type.slice(1)}</span>
                        <div class="move-details">Power: ${movePower}</div>
                    </button>
                `;
            }).join("");
            movesModal.classList.add("active");
            // Add listeners
            Array.from(movesGrid.querySelectorAll(".move-btn")).forEach(btn => {
                btn.onclick = (e) => {
                    const idx = parseInt(btn.getAttribute("data-move-idx"));
                    movesModal.classList.remove("active");
                    doMove("attack", idx);
                };
            });
        };
        defendBtn.onclick = () => doMove("defend");
        passBtn.onclick = () => doMove("pass");
        sacrificeBtn.onclick = () => doMove("sacrifice");
        surrenderBtn.onclick = () => {
            if (confirm("Are you sure you want to surrender?")) doMove("surrender");
        };
        cancelMoveBtn.onclick = () => movesModal.classList.remove("active");

        // Start the battle on load
        startBattle();
    }
});
