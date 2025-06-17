// Replace with your actual Railway backend URL
const BACKEND_URL = "poketactix-production.up.railway.app";

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
        const params = new URLSearchParams(window.location.search);
        const mode = params.get("mode") || "1v1";
        const player = params.get("player") || "";
        document.getElementById("arenaMode").textContent = `${mode.toUpperCase()} Battle Arena`;
        document.getElementById("arenaPlayerName").textContent = player;
        document.getElementById("arenaModeText").textContent = mode;
        document.getElementById("arenaPlayerText").textContent = player;
        document.getElementById("arenaFeature1").textContent = mode === "1v1" ? "Single Pokémon battles" : "Team-based 5v5 battles";
    }
});