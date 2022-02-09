import 'react-simple-keyboard/build/css/index.css';

import { Fireworks } from 'fireworks/lib/react';
import * as React from 'react';
import Keyboard from 'react-simple-keyboard';

import { makeGuessRequest, makeNewRandomGameRequest, makeStatusRequest } from '../api';
import { Attempt, MAX_ATTEMPTS, WORD_LENGTH } from '../types';
import { LetterGrid } from './LetterGrid';

type SimpleKeyboard = {
  addButtonTheme: (buttons: string, classes: string) => void;
  removeButtonTheme: (buttons: string, classes: string) => void;
};

enum GameState {
  InProgress,
  Won,
  Lost
}

export const Game = () => {
  const [input, setInput] = React.useState<string>('');
  const [attempts, setAttempts] = React.useState<Attempt[]>([]);
  const [correctWord, setCorrectWord] = React.useState<string | null>(null);
  const [gameState, setGameState] = React.useState<GameState>(GameState.InProgress);
  const keyboardRef = React.useRef<SimpleKeyboard | null>(null);

  React.useEffect(() => {
    (async () => {
      const [data] = await makeStatusRequest();
      if (data) {
        setAttempts(data.attempts);
        data.attempts.forEach(colorKeyboard);
      }
    })();
  }, []);

  const onKeyPress = async (key: string) => {
    if (key === '{bksp}' && input.length > 0) {
      setInput(input.substring(0, input.length - 1));
    } else if (key === '{enter}' && input.length === WORD_LENGTH) {
      const success = await makeGuess(input);
      success && setInput('');
    } else if (key.length === 1 && input.length < WORD_LENGTH) {
      setInput(`${input}${key}`);
    }
  };

  const handleUserKeyPress = (event: KeyboardEvent) => {
    // Needed for mac keyboard
    // TODO: test this doesn't cause double backspace on windows
    if (event.key === 'Backspace') {
      onKeyPress('{bksp}');
    }
  };

  React.useEffect(() => {
    window.addEventListener('keydown', handleUserKeyPress);
    return () => {
      window.removeEventListener('keydown', handleUserKeyPress);
    };
  }, [handleUserKeyPress]);

  const makeGuess = async (word: string): Promise<boolean> => {
    const [data, error] = await makeGuessRequest(word);
    if (error) {
      alert(error.message);
      return false;
    }
    if (data) {
      const hasWon = data.letters.every((l) => l.color === 'green');
      const hasLost = attempts.length === MAX_ATTEMPTS - 1 && !hasWon;
      setAttempts([...attempts, { word_guess: word, letters: data.letters }]);
      setCorrectWord(data.word ?? null);
      setTimeout(() => {
        hasWon && setGameState(GameState.Won);
        hasLost && setGameState(GameState.Lost);
        colorKeyboard(data);
      }, 2000);
    }
    return true;
  };

  const colorKeyboard = (attempt: Pick<Attempt, 'letters'>) => {
    attempt.letters.forEach((letter) => {
      keyboardRef.current?.addButtonTheme(letter.letter, `color-${letter.color}`);
    });
  };

  const onRandomClick = async (e: React.MouseEvent) => {
    e.preventDefault();
    const [data, error] = await makeNewRandomGameRequest();
    if (error) {
      alert(error.message);
      return;
    }
    if (data) {
      const currentUrl = window.location.href.split('?')[0];
      window.location.href = `${currentUrl}?token=${data.token}`;
    }
  };

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const resetGame = () => {
    setAttempts([]);
    keyboardRef.current?.removeButtonTheme(
      'q w e r t y u i o p a s d f g h j k l z x c v b n m',
      'color-grey color-yellow color-green'
    );
  };

  return (
    <>
      {gameState === GameState.Won && <Fireworks count={3} />}
      <LetterGrid attempts={attempts} input={input} />
      <Keyboard
        keyboardRef={(ref: SimpleKeyboard) => (keyboardRef.current = ref)}
        onChange={() => {}}
        onKeyPress={onKeyPress}
        physicalKeyboardHighlight
        physicalKeyboardHighlightPress
        layout={{
          default: ['q w e r t y u i o p', 'a s d f g h j k l', '{enter} z x c v b n m {bksp}']
        }}
        display={{
          '{enter}': 'ENTER',
          '{bksp}': '⬅'
        }}
      />
      {gameState === GameState.Lost && (
        <div className="correct-word">
          <span>Correct word was "{correctWord}"</span>
        </div>
      )}
      <div className="button-random">
        <a href="#" className="button" onClick={onRandomClick}>
          Random
        </a>
      </div>
    </>
  );
};
