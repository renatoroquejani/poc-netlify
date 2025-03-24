document.addEventListener('DOMContentLoaded', function() {
    // Exibir data e hora atual
    const now = new Date();
    document.getElementById('datetime').textContent = now.toLocaleString();
    
    // Adicionar evento ao botu00e3o
    const btn = document.getElementById('btn');
    btn.addEventListener('click', function() {
        alert('Deploy realizado com sucesso via pasta local!');
    });
});
