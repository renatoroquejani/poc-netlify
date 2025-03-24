document.addEventListener('DOMContentLoaded', function() {
    const deployForm = document.getElementById('deployForm');
    const usernameInput = document.getElementById('username');
    const previewSubdomain = document.getElementById('previewSubdomain');
    const deployButton = document.getElementById('deployButton');
    const deploySpinner = document.getElementById('deploySpinner');
    const resultCard = document.getElementById('resultCard');
    const deployResult = document.getElementById('deployResult');

    // Atualizar preview do subdomu00ednio quando o usuu00e1rio digitar
    usernameInput.addEventListener('input', function() {
        const username = this.value.trim() || 'usuario';
        previewSubdomain.textContent = `${username}.sites.kodestech.com.br`;
    });

    // Submeter formulu00e1rio
    deployForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        // Obter valores do formulu00e1rio
        const username = usernameInput.value.trim();
        const customDomain = document.getElementById('customDomain').value.trim();
        const s3Path = document.getElementById('s3Path').value.trim();
        
        // Validar campos obrigatu00f3rios
        if (!username) {
            alert('Por favor, informe o nome de usuu00e1rio.');
            usernameInput.focus();
            return;
        }
        
        if (!s3Path) {
            alert('Por favor, informe o caminho no S3.');
            document.getElementById('s3Path').focus();
            return;
        }
        
        // Preparar dados para envio
        const deployData = {
            username: username,
            custom_domain: customDomain,
            s3_path: s3Path
        };
        
        // Mostrar spinner e desabilitar botu00e3o
        deployButton.disabled = true;
        deploySpinner.classList.remove('d-none');
        resultCard.classList.add('d-none');
        
        // Enviar requisiu00e7u00e3o
        fetch('/api/deploy', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(deployData)
        })
        .then(response => response.json())
        .then(data => {
            // Exibir resultado
            showResult(data);
        })
        .catch(error => {
            // Exibir erro
            showError('Erro ao processar a requisiu00e7u00e3o: ' + error.message);
        })
        .finally(() => {
            // Esconder spinner e habilitar botu00e3o
            deployButton.disabled = false;
            deploySpinner.classList.add('d-none');
        });
    });

    // Funu00e7u00e3o para exibir resultado
    function showResult(data) {
        // Limpar resultado anterior
        deployResult.innerHTML = '';
        
        // Definir classe do caru00e1 de resultado
        const cardHeader = resultCard.querySelector('.card-header');
        if (data.success) {
            cardHeader.classList.remove('bg-danger');
            cardHeader.classList.add('bg-success');
            cardHeader.querySelector('h3').textContent = 'Deploy Iniciado com Sucesso';
        } else {
            cardHeader.classList.remove('bg-success');
            cardHeader.classList.add('bg-danger');
            cardHeader.querySelector('h3').textContent = 'Erro no Deploy';
        }
        
        // Construir HTML do resultado
        let resultHTML = `
            <div class="result-item">
                <span class="result-label">Status:</span> 
                <span class="${data.success ? 'result-success' : 'result-error'}">
                    ${data.success ? 'Sucesso' : 'Erro'}
                </span>
            </div>
            <div class="result-item">
                <span class="result-label">Mensagem:</span> 
                <span class="result-value">${data.message}</span>
            </div>
        `;
        
        // Adicionar informau00e7u00f5es adicionais se o deploy foi bem-sucedido
        if (data.success) {
            if (data.subdomain) {
                resultHTML += `
                    <div class="result-item">
                        <span class="result-label">Subdomu00ednio:</span> 
                        <span class="result-value">
                            <a href="https://${data.subdomain}" target="_blank">${data.subdomain}</a>
                        </span>
                    </div>
                `;
            }
            
            if (data.site_url) {
                resultHTML += `
                    <div class="result-item">
                        <span class="result-label">URL do Site:</span> 
                        <span class="result-value">
                            <a href="${data.site_url}" target="_blank">${data.site_url}</a>
                        </span>
                    </div>
                `;
            }
            
            if (data.custom_domain) {
                resultHTML += `
                    <div class="result-item">
                        <span class="result-label">Domu00ednio Personalizado:</span> 
                        <span class="result-value">
                            <a href="https://${data.custom_domain}" target="_blank">${data.custom_domain}</a>
                        </span>
                    </div>
                `;
            }
            
            if (data.deploy_id) {
                resultHTML += `
                    <div class="result-item">
                        <span class="result-label">ID do Deploy:</span> 
                        <span class="result-value">${data.deploy_id}</span>
                    </div>
                `;
            }
            
            resultHTML += `
                <div class="alert alert-info mt-3">
                    <strong>Nota:</strong> O deploy foi iniciado em segundo plano e pode levar alguns minutos para ser concluu00eddo.
                    Vocu00ea pode verificar o status do deploy no painel da Netlify.
                </div>
            `;
        }
        
        // Exibir resultado
        deployResult.innerHTML = resultHTML;
        resultCard.classList.remove('d-none');
        
        // Rolar para o resultado
        resultCard.scrollIntoView({ behavior: 'smooth' });
    }

    // Funu00e7u00e3o para exibir erro
    function showError(message) {
        const cardHeader = resultCard.querySelector('.card-header');
        cardHeader.classList.remove('bg-success');
        cardHeader.classList.add('bg-danger');
        cardHeader.querySelector('h3').textContent = 'Erro';
        
        deployResult.innerHTML = `
            <div class="result-item">
                <span class="result-label">Status:</span> 
                <span class="result-error">Erro</span>
            </div>
            <div class="result-item">
                <span class="result-label">Mensagem:</span> 
                <span class="result-value">${message}</span>
            </div>
        `;
        
        resultCard.classList.remove('d-none');
        resultCard.scrollIntoView({ behavior: 'smooth' });
    }
});
