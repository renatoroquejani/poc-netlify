function vTurbChangePlayer(){!vTurbOriginalPlayerIsMobile&&vTurbDeviceIsMobile||vTurbOriginalPlayerIsMobile&&!vTurbDeviceIsMobile?(vTurbPlayer=vTurbAlternativePlayer,vTurbSrcId=vTurbPlayer.id):vTurbPlayer=vTurbOriginalPlayer;var I=document.getElementById(`vid_${vTurbOriginalPlayer.id}`);I&&I.remove();var M=document.getElementById(`scr_${vTurbOriginalPlayer.id}`);M&&M.setAttribute("id",`scr_${vTurbSrcId}`)}function vTurbCreatSmartvdsElements(){var I,M,g;"1.7.9"===vTurbPlayer.version?document.getElementById(`vid_${vTurbPlayer.id}`)||(window,I=document,M=I.getElementById(`scr_${vTurbSrcId}`),(g=I.createElement("DIV")).id=`vid_${vTurbPlayer.id}`,M.parentElement.insertBefore(g,M)):(document.getElementById(`vid_${vTurbPlayer.id}`)||function(I,M,g){M=I.getElementById(`scr_${vTurbSrcId}`),(g=I.createElement("DIV")).id=`vid_${vTurbPlayer.id}`,g.style.position="relative",g.style.width="100%",g.style.padding=`${vTurbPlayer.video_aspect_ratio}% 0 0`,M.parentElement.insertBefore(g,M)}(document),document.getElementById(`thumb_${vTurbPlayer.id}`)||function(I,M,g){M=I.getElementById(`vid_${vTurbPlayer.id}`),(g=I.createElement("IMG")).id=`thumb_${vTurbPlayer.id}`,g.style.top="0",g.style.left="0",g.style.width="100%",g.style.height="100%",g.style.position="absolute",g.style.objectFit="cover",g.src=`https://images.converteai.net/${vTurbPlayer.thumbnail_key}`,M.appendChild(g)}(document),document.getElementById(`backdrop_${vTurbPlayer.id}`)||function(I,M,g){M=I.getElementById(`vid_${vTurbPlayer.id}`),(g=I.createElement("DIV")).id=`backdrop_${vTurbPlayer.id}`,g.style.top="0",g.style.left="0",g.style.width="100%",g.style.height="100%",g.style.position="absolute",g.style.backdropFilter="blur(5px)",g.style.webkitBackdropFilter="blur(5px)",M.appendChild(g)}(document))}function vTurbLoadSmrtvds(){var I,M,g,C;I=window,M=document,I.smrtvds||(g=I.smrtvds=function(){g.callMethod?g.callMethod.apply(g,arguments):g.queue.push(arguments)},I._smrtvds||(I._smrtvds=g),g.push=g,g.loaded=!0,g.version="1.1",g.queue=[],(C=M.createElement("script")).async=!0,C.src=`https://scripts.converteai.net/lib/js/smartplayer/${vTurbPlayer.version}/smartplayer.min.js`,M.getElementsByTagName("head")[0].appendChild(C)),window.smrtvds(`vid_${vTurbPlayer.id}`,vTurbPlayer.org_id,vTurbPlayer.video_id,vTurbPlayer.options)}function vTurbSmrtvds(){vTurbCreatSmartvdsElements(),vTurbLoadSmrtvds()}var vTurbOriginalPlayer={"id":"675c4ff26a0d47f470c80e1d","org_id":"554b8bee-f715-4c8a-85ef-f74cb5b80616","name":"WRB_VSL_FINAL_1.mp4","device_type":"desktop","video_aspect_ratio":"56.25","thumbnail_key":"554b8bee-f715-4c8a-85ef-f74cb5b80616/players/675c4ff26a0d47f470c80e1d/thumbnail.jpg","cover_key":"554b8bee-f715-4c8a-85ef-f74cb5b80616/players/675c4ff26a0d47f470c80e1d/cover.jpg","version":"v1","video_id":"675c4f360b9f9b72fcfb9fcd","options":{"autoplay":"smartplay","subtitle_active":!1,"smart_autoplay_template":"default","theme":"#ffa01c","foreground_color":"#ffffff","video":{"width":1920,"height":1080},"cdn":"cdn.converteai.net","conversion_params":["src"],"displays":{"big_play":!0,"play_pause":!1,"backward":!1,"subtitle_control":!1,"forward":!1,"volume":!0,"volume_bar":!0,"time":!1,"fullscreen":!0,"seekbar":!1,"seekbar_time":!0,"speed_control":!1},"callAction":[],"pixels":[],"thumbs":[],"headlines":[],"smart_autoplays":[{"id":"smart_autoplay_675c4ff26a0d47f470c80e1d_1","name":"Smart Autoplay 1","version":"2","number":1,"template":"default","background_color":"rgba(153, 204, 51, 0.57)","bottom_text":"Clique para ouvir","foreground_color":"#FFFFFF","top_text":"Seu v\xeddeo j\xe1 come\xe7ou","animation":{},"video_start_at":0,"video_end_at":10,"cover_key":"554b8bee-f715-4c8a-85ef-f74cb5b80616/players/675c4ff26a0d47f470c80e1d/smart_autoplay_675c4ff26a0d47f470c80e1d_1_cover.jpg","thumbnail_key":"554b8bee-f715-4c8a-85ef-f74cb5b80616/players/675c4ff26a0d47f470c80e1d/smart_autoplay_675c4ff26a0d47f470c80e1d_1_thumbnail.jpg","elements":[{"id":"smart_autoplay_675c4ff26a0d47f470c80e1d_1_element_0","height":480,"width":864,"x":528.0001467941329,"y":299.99966615142705,"order":1,"opacity":1,"rotation":0,"type":"box","properties":{"color":"rgba(255, 160, 28, 0.73)","radius":16,"border":{"size":4,"color":"rgb(255, 255, 255)","type":"solid"}}},{"id":"smart_autoplay_675c4ff26a0d47f470c80e1d_1_element_1","height":80,"width":864,"x":528.0001467941329,"y":342.6671830191262,"order":2,"opacity":1,"rotation":0,"type":"text","properties":{"size":53,"value":"Seu v\xeddeo j\xe1 come\xe7ou","color":"rgb(255, 255, 255)","weight":700,"align":"center"}},{"id":"smart_autoplay_675c4ff26a0d47f470c80e1d_1_element_2","height":80,"width":864,"x":528.0001467941329,"y":657.3328169808738,"order":3,"opacity":1,"rotation":0,"type":"text","properties":{"size":53,"value":"Clique para ouvir","color":"rgb(255, 255, 255)","weight":700,"align":"center"}},{"id":"smart_autoplay_675c4ff26a0d47f470c80e1d_1_element_3","height":192,"width":272,"x":824.0005071070046,"y":444.0002403709725,"order":4,"opacity":1,"rotation":0,"type":"image","properties":{"alt":"Seu v\xeddeo j\xe1 come\xe7ou","src":"data:image/svg+xml;base64,CiAgICAgIDxzdmcgdmVyc2lvbj0iMS4xIiBmaWxsPSIjRkZGRkZGIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIgogICAgICAgICAgeD0iMHB4IiB5PSIwcHgiIHdpZHRoPSI0Ni43NXB4IiBoZWlnaHQ9IjMyLjU2M3B4IiB2aWV3Qm94PSI3Ljk5OSA5LjA2MiA0Ni43NSAzMi41NjMiCiAgICAgICAgICBlbmFibGUtYmFja2dyb3VuZD0ibmV3IDcuOTk5IDkuMDYyIDQ2Ljc1IDMyLjU2MyIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIKICAgICAgPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgIEAtd2Via2l0LWtleWZyYW1lcyBCTElOSyB7CiAgICAgICAgICAgIDAlIHsgb3BhY2l0eTogMDsgfQogICAgICAgICAgICAzMyUgeyBvcGFjaXR5OiAxOyB9CiAgICAgICAgICAgIDY2JSB7IG9wYWNpdHk6IDE7IH0KICAgICAgICAgICAgMTAwJSB7IG9wYWNpdHk6IDA7IH0KICAgICAgICAgIH0KICAKICAgICAgICAgIEBrZXlmcmFtZXMgQkxJTksgewogICAgICAgICAgICAwJSB7IG9wYWNpdHk6IDA7IH0KICAgICAgICAgICAgMzMlIHsgb3BhY2l0eTogMTsgfQogICAgICAgICAgICA2NiUgeyBvcGFjaXR5OiAxOyB9CiAgICAgICAgICAgIDEwMCUgeyBvcGFjaXR5OiAwOyB9CiAgICAgICAgICB9CiAgCiAgICAgICAgICAuYW5pbWF0aW9uIC5ibGlua18xIHsKICAgICAgICAgICAgLXdlYmtpdC1hbmltYXRpb246IEJMSU5LIDJzIGluZmluaXRlOwogICAgICAgICAgICBhbmltYXRpb246IEJMSU5LIDJzIGluZmluaXRlOwogICAgICAgICAgICBvcGFjaXR5OiAwOwogICAgICAgICAgfQogIAogICAgICAgICAgLmFuaW1hdGlvbiAuYmxpbmtfMiB7CiAgICAgICAgICAgIC13ZWJraXQtYW5pbWF0aW9uOiBCTElOSyAycyBpbmZpbml0ZSAuM3M7CiAgICAgICAgICAgIGFuaW1hdGlvbjogQkxJTksgMnMgaW5maW5pdGUgLjNzOwogICAgICAgICAgICBvcGFjaXR5OiAwOwogICAgICAgICAgfQogIAogICAgICAgICAgLmFuaW1hdGlvbiAuYmxpbmtfMyB7CiAgICAgICAgICAgIC13ZWJraXQtYW5pbWF0aW9uOiBCTElOSyAycyBpbmZpbml0ZSAuNnM7CiAgICAgICAgICAgIGFuaW1hdGlvbjogQkxJTksgMnMgaW5maW5pdGUgLjZzOwogICAgICAgICAgICBvcGFjaXR5OiAwOwogICAgICAgICAgfQogIAogICAgICAgICAgLmFuaW1hdGlvbiAuc21hcnRwbGF5LXN2Zy1jb2xvciB7CiAgICAgICAgICAgIGZpbGw6ICcjRkZGRkZGJyAhaW1wb3J0YW50OwogICAgICAgICAgfQogIAogICAgICAgICAgLmFuaW1hdGlvbi5hZGp1c3RhYmxlIHsKICAgICAgICAgICAgYm9yZGVyOiA0cHggc29saWQgJyNGRkZGRkYnOwogICAgICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgCiAgICAgICAgPGcgY2xhc3M9ImFkanVzdGFibGUgZmcgYW5pbWF0aW9uIj4KICAgICAgICAgIDxwYXRoIGNsYXNzPSJzbWFydHBsYXktc3ZnLWNvbG9yIiBkPSJNNTMuMjQ5LDM5LjYxNmMtMC4xODYsMC0wLjM3MS0wLjA1MS0wLjUzNy0wLjE1N2wtNDMuNS0yNy43NWMtMC40NjYtMC4yOTctMC42MDMtMC45MTYtMC4zMDYtMS4zODFjMC4yOTgtMC40NjYsMC45MTctMC42MDEsMS4zODEtMC4zMDZsNDMuNSwyNy43NWMwLjQ2NywwLjI5NywwLjYwNCwwLjkxNiwwLjMwNywxLjM4MUM1My45MDEsMzkuNDUzLDUzLjU3OSwzOS42MTYsNTMuMjQ5LDM5LjYxNnoiPjwvcGF0aD4KICAgICAgICAgIDxwYXRoIGNsYXNzPSJibGlua18zIHNtYXJ0cGxheS1zdmctY29sb3IiIGQ9Ik00OC44OTYsMzMuNDY3bDEuNjk5LDEuMDg1YzMuNDk3LTcuNzkxLDIuMDczLTE3LjI3MS00LjMxMy0yMy42NTljLTAuMzkxLTAuMzkxLTEuMDIzLTAuMzkxLTEuNDE0LDBzLTAuMzkxLDEuMDIzLDAsMS40MTRDNTAuNTgxLDE4LjAxOSw1MS45MTMsMjYuNDYzLDQ4Ljg5NiwzMy40Njd6Ij48L3BhdGg+CiAgICAgICAgICA8cGF0aCBjbGFzcz0iYmxpbmtfMyBzbWFydHBsYXktc3ZnLWNvbG9yIiBkPSJNNDYuOTI2LDM2Ljk1NmMtMC42MTIsMC44NjMtMS4yODYsMS42OTUtMi4wNTksMi40NjljLTAuMzkyLDAuMzkxLTAuMzkyLDEuMDIzLDAsMS40MTRjMC4xOTQsMC4xOTUsMC40NSwwLjI5MywwLjcwNywwLjI5M2MwLjI1NiwwLDAuNTEyLTAuMDk4LDAuNzA2LTAuMjkzYzAuODc4LTAuODc4LDEuNjQyLTEuODI0LDIuMzMzLTIuODA3TDQ2LjkyNiwzNi45NTZ6Ij48L3BhdGg+CiAgICAgICAgICA8cGF0aCBjbGFzcz0iYmxpbmtfMiBzbWFydHBsYXktc3ZnLWNvbG9yIiBkPSJNNDIuNTQzLDI5LjQxNWwxLjc3NywxLjEzNWMxLjU0NS01LjMxNSwwLjIyOS0xMS4yOTMtMy45NTMtMTUuNDc2Yy0wLjM5Mi0wLjM5MS0xLjAyMy0wLjM5MS0xLjQxNCwwYy0wLjM5MiwwLjM5MS0wLjM5MiwxLjAyMywwLDEuNDE0QzQyLjQ1NCwxOS45ODcsNDMuNjM5LDI0LjkyNSw0Mi41NDMsMjkuNDE1eiI+PC9wYXRoPgogICAgICAgICAgPHBhdGggY2xhc3M9ImJsaW5rXzIgc21hcnRwbGF5LXN2Zy1jb2xvciIgZD0iTTQxLDMzLjE3NGMtMC41NjMsMC45NC0xLjIzNSwxLjgzNy0yLjA0NywyLjY0NmMtMC4zOTEsMC4zOTItMC4zOTEsMS4wMjMsMCwxLjQxNGMwLjE5NSwwLjE5NSwwLjQ1MSwwLjI5MywwLjcwNywwLjI5M3MwLjUxMi0wLjA5OCwwLjcwNy0wLjI5M2MwLjkxNi0wLjkxNCwxLjY3Ni0xLjkyNCwyLjMxNy0yLjk4NEw0MSwzMy4xNzR6Ij48L3BhdGg+CiAgICAgICAgICA8cGF0aCBjbGFzcz0iYmxpbmtfMSBzbWFydHBsYXktc3ZnLWNvbG9yIiBkPSJNMzUuNzcxLDI1LjA5NGwyLjAwMywxLjI3N2MwLjAxMi0wLjIwMywwLjAyOS0wLjQwNCwwLjAyOS0wLjYwOWMwLTMuMDc5LTEuMi01Ljk3NC0zLjM4MS04LjE1M2MtMC4zOTEtMC4zOTEtMS4wMjItMC4zOTEtMS40MTQsMGMtMC4zOTEsMC4zOTEtMC4zOTEsMS4wMjMsMCwxLjQxNEMzNC42NTIsMjAuNjY2LDM1LjYxMywyMi44MDIsMzUuNzcxLDI1LjA5NHoiPjwvcGF0aD4KICAgICAgICAgIDxwYXRoIGNsYXNzPSJibGlua18xIHNtYXJ0cGxheS1zdmctY29sb3IiIGQ9Ik0zNS4wODQsMjkuNDAxYy0wLjQ3NCwxLjE0NS0xLjE3MiwyLjE5Ny0yLjA3NiwzLjFjLTAuMzkxLDAuMzkxLTAuMzkxLDEuMDIzLDAsMS40MTRjMC4xOTUsMC4xOTUsMC40NTEsMC4yOTMsMC43MDcsMC4yOTNjMC4yNTcsMCwwLjUxMy0wLjA5OCwwLjcwNy0wLjI5M2MxLjAwOC0xLjAwNiwxLjc5NS0yLjE3LDIuMzYxLTMuNDNMMzUuMDg0LDI5LjQwMXoiPjwvcGF0aD4KICAgICAgICAgIDxwb2x5Z29uIGNsYXNzPSJzbWFydHBsYXktc3ZnLWNvbG9yIiBwb2ludHM9IjI4LjEyNCwyMC4yMTUgMjguMTI0LDE0Ljk5MSAyNC42MzUsMTcuOTkgICI+PC9wb2x5Z29uPgogICAgICAgICAgPHBhdGggY2xhc3M9InNtYXJ0cGxheS1zdmctY29sb3IiIGQ9Ik0yMC45MjEsMjAuMzY2aC02LjQyM2MtMC41NTMsMC0xLDAuNTA4LTEsMS4xMzV2OC4yMjljMCwwLjYyNywwLjQ0NywxLjEzNSwxLDEuMTM1aDcuMzc1bDYuMjUsNS44NzVWMjQuOTZMMjAuOTIxLDIwLjM2NnoiPjwvcGF0aD4KICAgICAgICA8L2c+CiAgICAgIDwvc3ZnPgogICAg"}}],"custom_player_preview":null}],"turbos":[],"smart_autoplay_elements":[{"id":"smart_autoplay_element_675c4ff26a0d47f470c80e1d_0","height":480,"width":864,"x":528.0001467941329,"y":299.99966615142705,"order":1,"opacity":1,"rotation":0,"type":"box","properties":{"color":"rgba(255, 160, 28, 0.73)","radius":16,"border":{"size":4,"color":"rgb(255, 255, 255)","type":"solid"}}},{"id":"smart_autoplay_element_675c4ff26a0d47f470c80e1d_1","height":80,"width":864,"x":528.0001467941329,"y":342.6671830191262,"order":2,"opacity":1,"rotation":0,"type":"text","properties":{"size":53,"value":"Seu v\xeddeo j\xe1 come\xe7ou","color":"rgb(255, 255, 255)","weight":700,"align":"center"}},{"id":"smart_autoplay_element_675c4ff26a0d47f470c80e1d_2","height":80,"width":864,"x":528.0001467941329,"y":657.3328169808738,"order":3,"opacity":1,"rotation":0,"type":"text","properties":{"size":53,"value":"Clique para ouvir","color":"rgb(255, 255, 255)","weight":700,"align":"center"}},{"id":"smart_autoplay_element_675c4ff26a0d47f470c80e1d_3","height":192,"width":272,"x":824.0005071070046,"y":444.0002403709725,"order":4,"opacity":1,"rotation":0,"type":"image","properties":{"alt":"Seu v\xeddeo j\xe1 come\xe7ou","src":"data:image/svg+xml;base64,CiAgICAgIDxzdmcgdmVyc2lvbj0iMS4xIiBmaWxsPSIjRkZGRkZGIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIgogICAgICAgICAgeD0iMHB4IiB5PSIwcHgiIHdpZHRoPSI0Ni43NXB4IiBoZWlnaHQ9IjMyLjU2M3B4IiB2aWV3Qm94PSI3Ljk5OSA5LjA2MiA0Ni43NSAzMi41NjMiCiAgICAgICAgICBlbmFibGUtYmFja2dyb3VuZD0ibmV3IDcuOTk5IDkuMDYyIDQ2Ljc1IDMyLjU2MyIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIKICAgICAgPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgIEAtd2Via2l0LWtleWZyYW1lcyBCTElOSyB7CiAgICAgICAgICAgIDAlIHsgb3BhY2l0eTogMDsgfQogICAgICAgICAgICAzMyUgeyBvcGFjaXR5OiAxOyB9CiAgICAgICAgICAgIDY2JSB7IG9wYWNpdHk6IDE7IH0KICAgICAgICAgICAgMTAwJSB7IG9wYWNpdHk6IDA7IH0KICAgICAgICAgIH0KICAKICAgICAgICAgIEBrZXlmcmFtZXMgQkxJTksgewogICAgICAgICAgICAwJSB7IG9wYWNpdHk6IDA7IH0KICAgICAgICAgICAgMzMlIHsgb3BhY2l0eTogMTsgfQogICAgICAgICAgICA2NiUgeyBvcGFjaXR5OiAxOyB9CiAgICAgICAgICAgIDEwMCUgeyBvcGFjaXR5OiAwOyB9CiAgICAgICAgICB9CiAgCiAgICAgICAgICAuYW5pbWF0aW9uIC5ibGlua18xIHsKICAgICAgICAgICAgLXdlYmtpdC1hbmltYXRpb246IEJMSU5LIDJzIGluZmluaXRlOwogICAgICAgICAgICBhbmltYXRpb246IEJMSU5LIDJzIGluZmluaXRlOwogICAgICAgICAgICBvcGFjaXR5OiAwOwogICAgICAgICAgfQogIAogICAgICAgICAgLmFuaW1hdGlvbiAuYmxpbmtfMiB7CiAgICAgICAgICAgIC13ZWJraXQtYW5pbWF0aW9uOiBCTElOSyAycyBpbmZpbml0ZSAuM3M7CiAgICAgICAgICAgIGFuaW1hdGlvbjogQkxJTksgMnMgaW5maW5pdGUgLjNzOwogICAgICAgICAgICBvcGFjaXR5OiAwOwogICAgICAgICAgfQogIAogICAgICAgICAgLmFuaW1hdGlvbiAuYmxpbmtfMyB7CiAgICAgICAgICAgIC13ZWJraXQtYW5pbWF0aW9uOiBCTElOSyAycyBpbmZpbml0ZSAuNnM7CiAgICAgICAgICAgIGFuaW1hdGlvbjogQkxJTksgMnMgaW5maW5pdGUgLjZzOwogICAgICAgICAgICBvcGFjaXR5OiAwOwogICAgICAgICAgfQogIAogICAgICAgICAgLmFuaW1hdGlvbiAuc21hcnRwbGF5LXN2Zy1jb2xvciB7CiAgICAgICAgICAgIGZpbGw6ICcjRkZGRkZGJyAhaW1wb3J0YW50OwogICAgICAgICAgfQogIAogICAgICAgICAgLmFuaW1hdGlvbi5hZGp1c3RhYmxlIHsKICAgICAgICAgICAgYm9yZGVyOiA0cHggc29saWQgJyNGRkZGRkYnOwogICAgICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgCiAgICAgICAgPGcgY2xhc3M9ImFkanVzdGFibGUgZmcgYW5pbWF0aW9uIj4KICAgICAgICAgIDxwYXRoIGNsYXNzPSJzbWFydHBsYXktc3ZnLWNvbG9yIiBkPSJNNTMuMjQ5LDM5LjYxNmMtMC4xODYsMC0wLjM3MS0wLjA1MS0wLjUzNy0wLjE1N2wtNDMuNS0yNy43NWMtMC40NjYtMC4yOTctMC42MDMtMC45MTYtMC4zMDYtMS4zODFjMC4yOTgtMC40NjYsMC45MTctMC42MDEsMS4zODEtMC4zMDZsNDMuNSwyNy43NWMwLjQ2NywwLjI5NywwLjYwNCwwLjkxNiwwLjMwNywxLjM4MUM1My45MDEsMzkuNDUzLDUzLjU3OSwzOS42MTYsNTMuMjQ5LDM5LjYxNnoiPjwvcGF0aD4KICAgICAgICAgIDxwYXRoIGNsYXNzPSJibGlua18zIHNtYXJ0cGxheS1zdmctY29sb3IiIGQ9Ik00OC44OTYsMzMuNDY3bDEuNjk5LDEuMDg1YzMuNDk3LTcuNzkxLDIuMDczLTE3LjI3MS00LjMxMy0yMy42NTljLTAuMzkxLTAuMzkxLTEuMDIzLTAuMzkxLTEuNDE0LDBzLTAuMzkxLDEuMDIzLDAsMS40MTRDNTAuNTgxLDE4LjAxOSw1MS45MTMsMjYuNDYzLDQ4Ljg5NiwzMy40Njd6Ij48L3BhdGg+CiAgICAgICAgICA8cGF0aCBjbGFzcz0iYmxpbmtfMyBzbWFydHBsYXktc3ZnLWNvbG9yIiBkPSJNNDYuOTI2LDM2Ljk1NmMtMC42MTIsMC44NjMtMS4yODYsMS42OTUtMi4wNTksMi40NjljLTAuMzkyLDAuMzkxLTAuMzkyLDEuMDIzLDAsMS40MTRjMC4xOTQsMC4xOTUsMC40NSwwLjI5MywwLjcwNywwLjI5M2MwLjI1NiwwLDAuNTEyLTAuMDk4LDAuNzA2LTAuMjkzYzAuODc4LTAuODc4LDEuNjQyLTEuODI0LDIuMzMzLTIuODA3TDQ2LjkyNiwzNi45NTZ6Ij48L3BhdGg+CiAgICAgICAgICA8cGF0aCBjbGFzcz0iYmxpbmtfMiBzbWFydHBsYXktc3ZnLWNvbG9yIiBkPSJNNDIuNTQzLDI5LjQxNWwxLjc3NywxLjEzNWMxLjU0NS01LjMxNSwwLjIyOS0xMS4yOTMtMy45NTMtMTUuNDc2Yy0wLjM5Mi0wLjM5MS0xLjAyMy0wLjM5MS0xLjQxNCwwYy0wLjM5MiwwLjM5MS0wLjM5MiwxLjAyMywwLDEuNDE0QzQyLjQ1NCwxOS45ODcsNDMuNjM5LDI0LjkyNSw0Mi41NDMsMjkuNDE1eiI+PC9wYXRoPgogICAgICAgICAgPHBhdGggY2xhc3M9ImJsaW5rXzIgc21hcnRwbGF5LXN2Zy1jb2xvciIgZD0iTTQxLDMzLjE3NGMtMC41NjMsMC45NC0xLjIzNSwxLjgzNy0yLjA0NywyLjY0NmMtMC4zOTEsMC4zOTItMC4zOTEsMS4wMjMsMCwxLjQxNGMwLjE5NSwwLjE5NSwwLjQ1MSwwLjI5MywwLjcwNywwLjI5M3MwLjUxMi0wLjA5OCwwLjcwNy0wLjI5M2MwLjkxNi0wLjkxNCwxLjY3Ni0xLjkyNCwyLjMxNy0yLjk4NEw0MSwzMy4xNzR6Ij48L3BhdGg+CiAgICAgICAgICA8cGF0aCBjbGFzcz0iYmxpbmtfMSBzbWFydHBsYXktc3ZnLWNvbG9yIiBkPSJNMzUuNzcxLDI1LjA5NGwyLjAwMywxLjI3N2MwLjAxMi0wLjIwMywwLjAyOS0wLjQwNCwwLjAyOS0wLjYwOWMwLTMuMDc5LTEuMi01Ljk3NC0zLjM4MS04LjE1M2MtMC4zOTEtMC4zOTEtMS4wMjItMC4zOTEtMS40MTQsMGMtMC4zOTEsMC4zOTEtMC4zOTEsMS4wMjMsMCwxLjQxNEMzNC42NTIsMjAuNjY2LDM1LjYxMywyMi44MDIsMzUuNzcxLDI1LjA5NHoiPjwvcGF0aD4KICAgICAgICAgIDxwYXRoIGNsYXNzPSJibGlua18xIHNtYXJ0cGxheS1zdmctY29sb3IiIGQ9Ik0zNS4wODQsMjkuNDAxYy0wLjQ3NCwxLjE0NS0xLjE3MiwyLjE5Ny0yLjA3NiwzLjFjLTAuMzkxLDAuMzkxLTAuMzkxLDEuMDIzLDAsMS40MTRjMC4xOTUsMC4xOTUsMC40NTEsMC4yOTMsMC43MDcsMC4yOTNjMC4yNTcsMCwwLjUxMy0wLjA5OCwwLjcwNy0wLjI5M2MxLjAwOC0xLjAwNiwxLjc5NS0yLjE3LDIuMzYxLTMuNDNMMzUuMDg0LDI5LjQwMXoiPjwvcGF0aD4KICAgICAgICAgIDxwb2x5Z29uIGNsYXNzPSJzbWFydHBsYXktc3ZnLWNvbG9yIiBwb2ludHM9IjI4LjEyNCwyMC4yMTUgMjguMTI0LDE0Ljk5MSAyNC42MzUsMTcuOTkgICI+PC9wb2x5Z29uPgogICAgICAgICAgPHBhdGggY2xhc3M9InNtYXJ0cGxheS1zdmctY29sb3IiIGQ9Ik0yMC45MjEsMjAuMzY2aC02LjQyM2MtMC41NTMsMC0xLDAuNTA4LTEsMS4xMzV2OC4yMjljMCwwLjYyNywwLjQ0NywxLjEzNSwxLDEuMTM1aDcuMzc1bDYuMjUsNS44NzVWMjQuOTZMMjAuOTIxLDIwLjM2NnoiPjwvcGF0aD4KICAgICAgICA8L2c+CiAgICAgIDwvc3ZnPgogICAg"}}],"mini_hooks":!1,"mini_hooks_elements":[],"resume":!0,"fake_bar":!0,"headline":!1,"turbo":!0,"turbo_speed":1.2,"turbo_auto_test":!1,"secure":!1,"smartplay_options":{"top_text":"Seu v\xeddeo j\xe1 come\xe7ou","bottom_text":"Clique para ouvir","foreground_color":"#FFFFFF","background_color":"rgba(153, 204, 51, 0.57)","start_at":0,"end_at":10,"animation":{},"custom_preview":null},"resume_options":{"play":"Continuar assistindo?","title":"Voc\xea j\xe1 come\xe7ou a assistir esse v\xeddeo","replay":"Assistir do in\xedcio?","disable_pause":!1,"foreground_color":"#FFFFFF","background_color":"#ffa01c"},"fake_bar_options":{"height":10,"alpha":2,"vbar_height":!0,"vbar_end":!0,"vbar_network":!0,"vbar_color":"#ffa01c"}}},vTurbSrcId="675c4ff26a0d47f470c80e1d",vTurbPlayer=vTurbOriginalPlayer,vTurbDeviceIsMobile=window.navigator.userAgent.match(/Mobile|iP(hone|od|ad)|Android|BlackBerry|IEMobile|Kindle|NetFront|Silk-Accelerated|(hpw|web)OS|Fennec|Minimo|Opera M(obi|ini)|Blazer|Dolfin|Dolphin|Skyfire|Zune/),vTurbOriginalPlayerIsMobile="mobile"===vTurbOriginalPlayer.device_type;vTurbDeviceIsMobile=vTurbDeviceIsMobile&&vTurbDeviceIsMobile[0],vTurbSmrtvds();