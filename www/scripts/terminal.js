  var term,
                    protocol,
                    socketURL,
                    socket,
                    charWidth,
                    charHeight,
                    cols,
                    rows,
                    szTimer,
                    sessionStats;

                var terminalContainer = document.getElementById('terminal-container');
                createTerminal();
                
                                console.log('Test3');

                function createTerminal() {
                  //while (terminalContainer.children.length) {
                  //  terminalContainer.removeChild(terminalContainer.children[0]);
                  //}
                  console.log('Test4');
                  
                  term = new Terminal({
                    cursorBlink: true,
                    screenKeys: true,
                    /* Initial terminal size to match the session originator's terminal */
                    cols: 100,
                    rows: 30,
                  });
                  protocol = (location.protocol === 'https:') ? 'wss://' : 'ws://';
                  sessionId = location.pathname.split("/").pop()
                  socketURL = protocol + location.hostname + ((location.port) ? (':' + location.port) : '') + 
                    '/shell/' + '?' + document.cookie;
                    
                    console.log('Test');
                    console.log(socketURL);

                  socket = new WebSocket(socketURL);

                  term.open(terminalContainer);
                  term.fit();

                  socket.onopen  = runRealTerminal;
                  socket.onclose = runRealTerminal;
                  socket.onerror = runRealTerminal;

                  var initialGeometry = term.proposeGeometry()
                  cols = initialGeometry.cols;
                  rows = initialGeometry.rows;

                  charWidth = Math.ceil(term.element.offsetWidth / cols);
                  charHeight = Math.ceil(term.element.offsetHeight / rows);

                  //szTimer = window.setInterval(updateTerminalSize, 1000)
                }
                
                
                
                function runRealTerminal() {
                  term.attach(socket);
                  term._initialized = true;
                  setTerminalSize(cols, rows);
                }

                function setTerminalSize (cols, rows) {
                  var width = (cols * charWidth).toString() + 'px',
                      height = (rows * charHeight).toString() + 'px';

                  terminalContainer.style.width = width;
                  terminalContainer.style.height = height;
                  term.resize(cols, rows);
                }

                // gets called once a second to re-size the terminal to match it to
                // the session originator's
                
                
