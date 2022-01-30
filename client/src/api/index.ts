import { Letter } from '../types';

const URL = 'http://localhost:8080/api';

type ErrorResponse = {
  error_code: string;
  message: string;
};

type GuessSuccessResponse = {
  token: string;
  letters: [Letter, Letter, Letter, Letter, Letter];
};

type GuessResponse = ErrorResponse | GuessSuccessResponse;

const responseIsError = (response: GuessResponse): response is ErrorResponse => {
  return 'error_code' in response;
};

const getToken = () => {
  const params = new URLSearchParams(window.location.search);
  const token = params.get('token');
  if (!token) {
    throw new Error('No token');
  }
  return token;
};

export const makeGuessRequest = async (
  word: string
): Promise<[GuessSuccessResponse, null] | [null, ErrorResponse]> => {
  const response = await fetch(`${URL}/game/guess/${word}`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Authorization: `Bearer ${getToken()}`
    }
  });
  const data = (await response.json()) as GuessResponse;
  if (responseIsError(data)) {
    return [null, data];
  }
  return [data, null];
};

export const makeStartRequest = async (): Promise<void> => {
  await fetch(`${URL}/game/start`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Authorization: `Bearer ${getToken()}`
    }
  });
};
