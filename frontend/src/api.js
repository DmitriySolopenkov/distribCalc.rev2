import axios from 'axios';

const API_URL = `http://${process.env.REACT_APP_API_SERVER}/api/v1/`;

class ApiService {
  getTasks() {
    return axios.get(API_URL + "task", { headers: this.authHeader() });
  }
  addTask(expression) {
    return axios.post(API_URL + "task", {expression: expression}, { headers: this.authHeader() })
  }

  getAgents() {
    return axios.get(API_URL + "agent");
  }

  login(login, password) {
    return axios.post(API_URL + "login", {
      login,
      password
    }).then(response => {
      return response.data;
    });
  }

  register(login, password) {
    return axios.post(API_URL + "register", {
      login,
      password
    }).then(response => {
      return response.data;
    });
  }

  getCurrentUser() {
    return JSON.parse(localStorage.getItem('user'));
  }

  authHeader() {
    const user = JSON.parse(localStorage.getItem('user'));

    if (user && user.token) {
      return {
        ContentType: 'application/json',
        Accept: 'application/json',
        Authorization: user.token
      };
    } else {
      return {};
    }
  }
}

const apiService = new ApiService();

export default apiService;
