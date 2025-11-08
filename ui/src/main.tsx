import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

import { RouterProvider } from "react-router/dom";
import {createBrowserRouter} from "react-router";
import LandingPage from "./pages/LandingPage.tsx";

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
