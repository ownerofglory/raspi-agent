import { createContext, type Dispatch, type SetStateAction } from "react";

/**
 * Represents the authenticated user's data.
 */
interface Auth {
    /** The user's authentication token */
    token: string;

    /** The user's unique identifier */
    id: string;
}

/**
 * Shape of the Auth context value.
 */
export interface AuthContextProps {
    /** Current authentication state */
    auth: Auth | undefined;

    /** Setter for authentication state */
    setAuth: Dispatch<SetStateAction<Auth | undefined>>;
}

/**
 * Authentication Context
 *
 * Provides `auth` and `setAuth` to children via a Context Provider.
 * Default value is `null` until wrapped with an AuthProvider.
 */
const AuthCtx = createContext<AuthContextProps | null>(null);

export default AuthCtx;
