import './App.css';

import { useSearchParams } from 'react-router-dom';

import { CreateGame } from './components/CreateGame';
import { Game } from './components/Game';

function App() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');
  return (
    <div className="App">
      {token && <Game />}
      {!token && <CreateGame />}
    </div>
  );
}

export default App;
