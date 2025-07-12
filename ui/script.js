function login(event) {
  event.preventDefault();
  fetch('/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: document.getElementById('email').value,
      password: document.getElementById('password').value
    })
  })
    .then(res => res.json())
    .then(data => {
      alert(data)
      console.log(JSON.stringify(data))
      if (data.data.access_token) {
        localStorage.setItem('access_token', data.data.access_token);
        localStorage.setItem('refresh_token', data.data.refresh_token);
        alert("Login successful!");
      } else {
        alert("Login failed");
      }
    })
    .catch(error => {
      console.error('Login error:', error);
      alert("Login failed");
    });
}

function signup(event) {
  event.preventDefault();
  fetch('/signup', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: document.getElementById('email').value,
      password: document.getElementById('password').value
    })
  })
    .then(res => res.text())
    .then(msg => alert(msg));
}

async function fetchProfile() {
  let token = localStorage.getItem('access_token');

  try {
    let response = await fetch('/me', {
      headers: { 'Authorization': 'Bearer ' + token }
    });

    if (response.status === 401) {
      // Token might be expired, try to refresh
      const refreshed = await refreshToken();
      if (refreshed) {
        token = localStorage.getItem('access_token');
        response = await fetch('/me', {
          headers: { 'Authorization': 'Bearer ' + token }
        });
      } else {
        throw new Error('Authentication failed');
      }
    }

    const data = await response.json();
    document.getElementById('profile').innerText = JSON.stringify(data, null, 2);
  } catch (error) {
    console.error('Profile fetch error:', error);
    alert('Please login again');
  }
}

async function refreshToken() {
  const refreshToken = localStorage.getItem('refresh_token');

  if (!refreshToken) {
    return false;
  }

  try {
    const response = await fetch('/refresh-token', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        refresh_token: refreshToken
      })
    });

    if (response.ok) {
      const data = await response.json();
      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('refresh_token', data.refresh_token);
      return true;
    } else {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      return false;
    }
  } catch (error) {
    console.error('Token refresh error:', error);
    return false;
  }
}