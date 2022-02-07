import './CreateGame.css';

import * as React from 'react';

import { makeNewGameRequest } from '../api';

export const CreateGame = () => {
  const [input, setInput] = React.useState<string>('');
  const [error, setError] = React.useState<string | null>(null);
  const [token, setToken] = React.useState<string | null>(null);
  const onSubmit = async () => {
    if (input.length !== 5) {
      setError('Word must be 5 letters long');
    }
    const [response, error] = await makeNewGameRequest(input);
    if (error) {
      setError(error.message);
    }
    if (response) {
      setToken(response.token);
      setError(null);
      setInput('');
    }
  };
  const onChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setError(null);
    setInput(event.target.value);
  };
  const onKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter') {
      onSubmit();
    }
  };
  return (
    <>
      <h1>Create a game</h1>
      <div>Pick a 5 letter word to send to a friend</div>
      <input
        className="word-input"
        type="text"
        maxLength={5}
        minLength={5}
        onChange={onChange}
        onKeyDown={onKeyDown}
        value={input}
      />
      <input className="word-submit" type="submit" onClick={onSubmit} />
      {error && <div className="word-error">{error}</div>}
      {token && (
        <div>
          Copy this link and send it to your friend:
          <textarea
            rows={5}
            value={`${window.location.origin}/?token=${token}`}
            className="token-input"
          ></textarea>
        </div>
      )}
    </>
  );
};
