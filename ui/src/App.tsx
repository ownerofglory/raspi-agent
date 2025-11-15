import './App.css'
import { Outlet } from "react-router";
import AuthCtx from "./context/auth";
import {useState} from "react";

/**
 * Root application component.
 *
 * Provides the Auth context to the rest of the app.
 */
function App() {
    /**
     * Authentication state.
     *
     * `auth` contains the user's token and id when logged in,
     * or `undefined` when not authenticated.
     */
    const [auth, setAuth] = useState<{
        token: string;
        id: string;
    } | undefined>(undefined);

    return (
        <AuthCtx.Provider value={{ auth, setAuth }}>
            <Outlet />
        </AuthCtx.Provider>
    );
}

export default App;
