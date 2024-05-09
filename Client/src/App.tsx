import { BrowserRouter, Routes } from 'react-router-dom';
import './App.css';
import Theme from './Theme';
import Footer from './components/Shared/Footer/Footer';
import Header from './components/Shared/Header/Header';

function App() {
  return (
    <Theme>
      <div className="App">
        <BrowserRouter>
          <Header />
          <br />
          <Routes>
          </Routes>
          <br />
          <Footer />
        </BrowserRouter>
      </div>
    </Theme>
  );
}

export default App;
