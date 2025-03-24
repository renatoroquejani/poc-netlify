
// document.addEventListener('click', function(event) {
//     let element = event.target;
//     while (element) {
//         if (element.getAttribute('data-framer-name') === 'container-video-vturb' && !element.classList.contains('video-ativado')) {
//             event.preventDefault();
//             event.stopPropagation();
//             element.innerHTML = '<div id=\"onm-video\"><div id=\"vid_660d93bb9074c30008a5d3b0\" style=\"position:relative;width:100%;padding: 56.25% 0 0;\"> <img id=\"thumb_660d93bb9074c30008a5d3b0\" src=\"https://images.converteai.net/554b8bee-f715-4c8a-85ef-f74cb5b80616/players/660d93bb9074c30008a5d3b0/thumbnail.jpg\" style=\"position:absolute;top:0;left:0;width:100%;height:100%;object-fit:cover;display:block;\"> <div id=\"backdrop_660d93bb9074c30008a5d3b0\" style=\"position:absolute;top:0;width:100%;height:100%;-webkit-backdrop-filter:blur(5px);backdrop-filter:blur(5px);\"></div> </div> ' +
//                 '<script type=\"text/javascript\" id=\"scr_660d93bb9074c30008a5d3b0\">  </script>' +
//                 '</div>';
//             var s=document.createElement("script"); s.src="https://scripts.converteai.net/554b8bee-f715-4c8a-85ef-f74cb5b80616/players/660d93bb9074c30008a5d3b0/player.js", s.async=!0,document.head.appendChild(s);
//             element.classList.add('video-ativado');
//             break;
//         }
//         element = element.parentElement;
//     }
// });

console.log("Definindo variáveis");

/* ALTERE O VALOR 10 PARA OS SEGUNDOS EM QUE AS SEÇÕES VÃO APARECER */
let SECONDS_TO_DISPLAY = 10;
let CLASS_TO_DISPLAY = '.esconder';
/* DAQUI PARA BAIXO NAO PRECISA ALTERAR */
let attempts = 0;
let elsHiddenList = [];
let elsDisplayed = false;
let elsHidden = document.querySelectorAll(CLASS_TO_DISPLAY);
let alreadyDisplayedKey = 'alreadyElsDisplayed' + SECONDS_TO_DISPLAY;
let alreadyElsDisplayed = null;


const onmHidden = document.querySelectorAll('[data-framer-name="esconder"]');
onmHidden.forEach((element) => {
    element.classList.add('esconder');
});

setTimeout(function(){

    const seletorIframe = document.querySelector('[data-framer-name="container-vturb"] iframe');
    // console.log(seletorIframe);

        console.log("Tudo carregado");

            try {
                alreadyElsDisplayed = localStorage.getItem(alreadyDisplayedKey);
            } catch (e) {
                console.warn('Failed to read data from localStorage: ', e);
            }
            setTimeout(function () {
                elsHiddenList = Array.prototype.slice.call(elsHidden);
            }, 0);
            var showHiddenElements = function () {
                elsDisplayed = true;
                onmHidden.forEach((e) => (e.style.display = 'flex'));
                try {
                    localStorage.setItem(alreadyDisplayedKey, true);
                } catch (e) {
                    console.warn('Failed to save data in localStorage: ', e);
                }
            };
            var startWatchVideoProgress = function () {
                if (typeof seletorIframe.contentWindow.smartplayer === 'undefined' || !(seletorIframe.contentWindow.smartplayer.instances && seletorIframe.contentWindow.smartplayer.instances.length)) {
                    console.log('entrei no if typeof smartplayer');
                    if (attempts >= 10) return;
                    attempts += 1;
                    return setTimeout(function () {
                        startWatchVideoProgress();
                    }, 1000);
                }
                seletorIframe.contentWindow.smartplayer.instances[0].on('timeupdate', () => {
                    console.log('entrei no Smartplayer Instances');
                    if (elsDisplayed || seletorIframe.contentWindow.smartplayer.instances[0].smartAutoPlay) return;
                    if (seletorIframe.contentWindow.smartplayer.instances[0].video.currentTime < SECONDS_TO_DISPLAY) return;
                    console.log('Exibindo');
                    showHiddenElements();
                });
            };
            if (alreadyElsDisplayed === 'true') {
                setTimeout(function () {
                    showHiddenElements();
                }, 100);
            } else {
                startWatchVideoProgress();
            }


            console.log("Tudo certo.");

}, 3000);
