/**
 * MangaHub Premium Dashboard JS
 * Pure Vanilla JS - No Frameworks allowed!
 */

let token = null;
let ws = null;

const LEDS = {
    HTTP: document.getElementById('led-http'),
    TCP: document.getElementById('led-tcp'),
    UDP: document.getElementById('led-udp'),
    WS: document.getElementById('led-ws'),
    GRPC: document.getElementById('led-grpc')
};

function glow(protocol) {
    if (LEDS[protocol]) {
        LEDS[protocol].classList.add('led-glow');
        setTimeout(() => {
            LEDS[protocol].classList.remove('led-glow');
        }, 500);
    }
}

// 🔐 Authentication
document.getElementById('login-btn').addEventListener('click', async () => {
    const user = document.getElementById('username').value;
    const pass = document.getElementById('password').value;

    try {
        const resp = await fetch('/api/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: user, password: pass })
        });

        if (resp.ok) {
            const data = await resp.json();
            token = data.token;
            glow('HTTP');
            document.getElementById('login-overlay').style.display = 'none';
            initWebSocket();
            showToast('Welcome Back, Architect. 🌸');
        } else {
            alert('Access Denied. Incorrect credentials.');
        }
    } catch (err) {
        console.error('Login failed', err);
    }
});

// 📡 Real-time Updates (WebSockets)
function initWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/api/chat?token=${token}&manga_id=1`);

    ws.onopen = () => {
        glow('WS');
        document.getElementById('feed-status').textContent = 'Live (WebSocket Connected)';
        document.getElementById('feed-status').style.color = '#FBCFE8';
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        
        // Check if it's a chat message or a system event
        if (data.username && data.content) {
            appendChat(data);
            glow('WS');
        } else if (data.manga_id) {
            appendEvent(data);
            glow('UDP'); // Treat manga.new broadcasts as "UDP/Bus" events
            showToast(`🚀 New Release: ${data.message || 'Manga Updated'}`);
        }
    };

    ws.onclose = () => {
        document.getElementById('feed-status').textContent = 'Disconnected';
        document.getElementById('feed-status').style.color = '#ff4444';
        setTimeout(initWebSocket, 3000);
    };
}

function appendChat(msg) {
    const container = document.getElementById('chat-messages');
    const div = document.createElement('div');
    const isMe = msg.username === 'admin' || msg.username === document.getElementById('username').value;
    
    div.className = `flex flex-col ${isMe ? 'items-end' : 'items-start'}`;
    div.innerHTML = `
        <span class="text-[10px] text-slate-500 uppercase font-bold px-2">${msg.username}</span>
        <div class="px-4 py-2 rounded-lg text-sm ${isMe ? 'bg-pink-pastel text-subtle-black font-medium' : 'bg-slate-800 text-slate-300'}">
            ${msg.content}
        </div>
    `;
    container.appendChild(div);
    container.scrollTop = container.scrollHeight;
}

function appendEvent(evt) {
    const container = document.getElementById('event-feed');
    const div = document.createElement('div');
    div.className = 'text-xs border-l-2 border-pink-pastel pl-3 py-1 bg-white/5 rounded-r-md animate-pulse';
    div.innerHTML = `
        <span class="text-pink-pastel font-bold uppercase mr-2">[RELEASE]</span>
        <span class="text-slate-300">${evt.message}</span>
        <div class="text-[10px] text-slate-500 mt-1">${new Date().toLocaleTimeString()}</div>
    `;
    container.prepend(div);
}

// 💬 Chat Actions
document.getElementById('send-chat-btn').addEventListener('click', sendChat);
document.getElementById('chat-input').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') sendChat();
});

function sendChat() {
    const input = document.getElementById('chat-input');
    const content = input.value.trim();
    if (!content || !ws) return;

    ws.send(content);
    input.value = '';
}

// 🔨 Create Manga (Broadcast)
document.getElementById('create-manga-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const title = document.getElementById('manga-title').value;
    const author = document.getElementById('manga-author').value;
    const description = document.getElementById('manga-desc').value;

    try {
        const resp = await fetch('/api/manga', {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ title, author, description })
        });

        if (resp.ok) {
            glow('HTTP');
            glow('GRPC'); // Usually gRPC handles this internally
            showToast(`Manga "${title}" created successfully! 🌸`);
            document.getElementById('create-manga-form').reset();
        }
    } catch (err) {
        showToast('Failed to create manga. Check connection.');
    }
});

function showToast(msg) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = 'pink-toast';
    toast.textContent = msg;
    container.appendChild(toast);
    
    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transform = 'translateX(100%)';
        toast.style.transition = 'all 0.5s ease';
        setTimeout(() => toast.remove(), 500);
    }, 4000);
}
