import './App.css';
import 'react-simple-keyboard/build/css/index.css';

import * as React from 'react';
import Keyboard from 'react-simple-keyboard';

import { makeGuessRequest } from './api';
import { LetterGrid } from './components/LetterGrid';
import { Attempt, WORD_LENGTH } from './types';

function App() {
  const [input, setInput] = React.useState<string>('');
  const [attempts, setAttempts] = React.useState<Attempt[]>([]);

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
      return false;
    }
    if (data) {
      const hasWon = data.letters.every((l) => l.color === 'green');
      if (hasWon) {
        alert('you won');
        resetGame();
      } else {
        setAttempts([...attempts, { word, letters: data.letters }]);
      }
    }
    return true;
  };

  const resetGame = () => {
    // TODO: show emoji stats
    setAttempts([]);
  };

  return (
    <div className="App">
      <LetterGrid attempts={attempts} input={input} />
      <Keyboard
        onChange={() => {}}
        onKeyPress={onKeyPress}
        physicalKeyboardHighlight
        physicalKeyboardHighlightPress
        layout={{
          default: [
            'q w e r t y u i o p',
            'a s d f g h j k l',
            '{enter} z x c v b n m {bksp} {delete}'
          ]
        }}
      />
    </div>
  );
}

export default App;
