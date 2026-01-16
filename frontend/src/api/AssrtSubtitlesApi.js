import { createRequest } from 'src/utils/http';
import config from 'src/config';

class AssrtSubtitlesApi {
  BaseUrl = config.BACKEND_URL;

  http(url, ...option) {
    return createRequest(`${this.BaseUrl}/v1${url}`, ...option);
  }

  search = (data) => this.http('/subtitles/assrt/search', data, 'POST');

  download = (uid) => this.http('/subtitles/assrt/download', { uid }, 'GET', { responseType: 'blob' });
}

export default new AssrtSubtitlesApi();
