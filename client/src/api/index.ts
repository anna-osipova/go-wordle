import { Letter } from '../types';

const URL = 'http://localhost:8080/api';
const token =
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ3b3JkIjoicm9hc3QiLCJhdHRlbXB0cyI6MCwiZXhwIjoxNjQzMDcwNjU3LCJpYXQiOjE2NDI4OTc4NTcsImlzcyI6IkFubmEifQ.Qpg53BcSCWZ7Uag7EGNEXlKlgP8WtTP550W9zUnppT4';

type ErrorResponse = {
  error: string;
};

type GuessSuccessResponse = {
  token: string;
  letters: [Letter, Letter, Letter, Letter, Letter];
};

type GuessResponse = ErrorResponse | GuessSuccessResponse;

const responseIsError = (response: GuessResponse): response is ErrorResponse => {
  return 'error' in response;
};

export const makeGuessRequest = async (
  word: string
): Promise<[GuessSuccessResponse, null] | [null, ErrorResponse]> => {
  const response = await fetch(`${URL}/game/guess/${word}?token=${token}`, { method: 'POST' });
  const data = (await response.json()) as GuessResponse;
  if (responseIsError(data)) {
    return [null, data];
  }
  return [data, null];
};
