import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';
import Theme from './Theme';
import Footer from './components/Shared/Footer/Footer';
import Header from './components/Shared/Header/Header';
import NoPage from './pages/NoPage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import { Toaster } from 'react-hot-toast';
import HomePage from './pages/HomePage';
import MupPage from './pages/MupPage';
import PolicePage from './pages/PolicePage';
import CourtPage from './pages/CourtPage';
import StatisticsPage from './pages/StatisticsPage';
import PasswordResetPage from './pages/PasswordResetPage';

function App() {
  return (
    <Theme>
      <div className="App">
        <BrowserRouter>
          <Header />
          <br />
          <Toaster position="top-right" />
          <Routes>
            <Route index element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route path='/reset-password' element={<PasswordResetPage />} />
            <Route path="/home" element={<HomePage />} />
            <Route path="/home/mup" element={<MupPage />} />
            <Route path="/home/police" element={<PolicePage />} />
            <Route path="/home/court" element={<CourtPage />} />
            <Route path="/home/statistics" element={<StatisticsPage />} />
            <Route path="*" element={<NoPage />} />
          </Routes>
          <br />
          <Footer />
        </BrowserRouter>
      </div>
    </Theme>
  );
}

export default App;
