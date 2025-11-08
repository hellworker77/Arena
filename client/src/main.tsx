import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import "./index.css";
import "./i18n/i18n.js";
import App from './components/pages/app/app.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
        <App />
  </StrictMode>,
)
