import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

import { RouterProvider } from "react-router/dom";
import {createBrowserRouter} from "react-router";
import LandingPage from "./pages/LandingPage.tsx";
import LoginPage from "./pages/auth/LoginPage.tsx";
import Dashboard from "./pages/dashboard/DashboardPage.tsx";

const router = createBrowserRouter(
    [
        {
            path: "/",
            element: <App />,
            children: [
                {
                    index: true,
                    element: <LandingPage />,
                },
            ],
        },
        {
            path: "/auth/login",
            element: <LoginPage/>
        },
        {
            path: "/devices",
            element: <Dashboard/>
        }
    ],
    {
        basename: "/raspi-agent/ui",
    }
);


createRoot(document.getElementById('root')!).render(
  <StrictMode>
      <RouterProvider router={router} />
  </StrictMode>,
)
