const refreshToken = async () => {
  const refreshToken = localStorage.getItem('refresh_token');
  if (!refreshToken) {
    return false;
  }
  try {
    const response = await fetch('/refresh-token', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken })
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
  } catch (err) {
    console.error('Token refresh error:', err);
    return false;
  }
};

const createLoginApp = () => ({
  data() {
    return { email: '', password: '' };
  },
  methods: {
    async login() {
      try {
        const res = await fetch('/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email: this.email, password: this.password })
        });
        const data = await res.json();
        if (data.data && data.data.access_token) {
          localStorage.setItem('access_token', data.data.access_token);
          localStorage.setItem('refresh_token', data.data.refresh_token);
          alert('Login successful!');
        } else {
          alert('Login failed');
        }
      } catch (err) {
        console.error('Login error:', err);
        alert('Login failed');
      }
    }
  }
});

const createSignupApp = () => ({
  data() {
    return { email: '', password: '' };
  },
  methods: {
    async signup() {
      try {
        const res = await fetch('/signup', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email: this.email, password: this.password })
        });
        const msg = await res.text();
        alert(msg);
      } catch (err) {
        console.error('Signup error:', err);
        alert('Signup failed');
      }
    }
  }
});

const createProfileApp = () => ({
  data() {
    return { profile: '' };
  },
  methods: {
    async fetchProfile() {
      let token = localStorage.getItem('access_token');
      try {
        let response = await fetch('/me', {
          headers: { 'Authorization': 'Bearer ' + token }
        });
        if (response.status === 401) {
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
        this.profile = JSON.stringify(data, null, 2);
      } catch (err) {
        console.error('Profile fetch error:', err);
        alert('Please login again');
      }
    },
    logout() {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      this.profile = '';
      alert('Logged out');
    }
  }
});

// Mount apps conditionally based on element presence
window.addEventListener('DOMContentLoaded', () => {
  if (document.getElementById('login-app')) {
    Vue.createApp(createLoginApp()).mount('#login-app');
  }
  if (document.getElementById('signup-app')) {
    Vue.createApp(createSignupApp()).mount('#signup-app');
  }
  if (document.getElementById('profile-app')) {
    Vue.createApp(createProfileApp()).mount('#profile-app');
  }
});

function logout() {
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
  alert('Logged out');
}
