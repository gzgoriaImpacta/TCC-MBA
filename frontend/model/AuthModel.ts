export class AuthModel {
    accerToken: string;
    refreshToken: string;

    constructor(
        accerToken: string, 
        refreshToken: string
    ) {
        this.accerToken = accerToken;
        this.refreshToken = refreshToken;
    }
}