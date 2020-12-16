import axios, { AxiosInstance } from "axios";
import config from "../config";

const instance: AxiosInstance = axios.create({
  baseURL: `${config.ApiEndpoint}${config.ApiPrefix}`,
  timeout: 1000,
  headers: {}
});

export default instance;
