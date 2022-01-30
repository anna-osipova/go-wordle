import './App.css';
import 'react-simple-keyboard/build/css/index.css';

import * as React from 'react';
import Keyboard from 'react-simple-keyboard';

import { makeGuessRequest, makeStartRequest } from './api';
import { LetterGrid } from './components/LetterGrid';
import { Attempt, WORD_LENGTH } from './types';

type SimpleKeyboard = {
  addButtonTheme: (buttons: string, classes: string) => void;
  removeButtonTheme: (buttons: string, classes: string) => void;
};

function App() {
  const [input, setInput] = React.useState<string>('');
  const [attempts, setAttempts] = React.useState<Attempt[]>([]);
  const keyboardRef = React.useRef<SimpleKeyboard | null>(null);

  React.useEffect(() => {
    makeStartRequest();
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
      setAttempts([...attempts, { word, letters: data.letters }]);
      const hasWon = data.letters.every((l) => l.color === 'green');
      if (hasWon) {
        alert('you won');
        resetGame();
      } else {
        data.letters.forEach((letter) => {
          keyboardRef.current?.addButtonTheme(letter.letter, `color-${letter.color}`);
        });
      }
    }
    return true;
  };

  const resetGame = () => {
    // TODO: show emoji stats
    setAttempts([]);
    keyboardRef.current?.removeButtonTheme(
      'q w e r t y u i o p a s d f g h j k l z x c v b n m',
      'color-grey color-yellow color-green'
    );
  };

  return (
    <div className="App">
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
      />
    </div>
  );
}

export default App;
