

export const selectUser = (state: NetAlertState) => state.NetworkInfo;
export const selectProfiles = (state: NetAlertState) => state.Profiles;

export interface NetAlertState{
  isLoading: boolean;
  error: any;
  NetworkInfo?: NetworkInfo;
  Profiles?: Profile[];
}

export interface NetworkCard {
  Name: string;
  Mac: string;
}

export interface NetworkInfo {
  SSID: string;
  BSSID: string;
  GatewayIP: string;
  GatewayMAC: string;
  NetworkCards: NetworkCard[];
}



export interface Profile {
  Mac: string;
  NickName: string;
  CreateDate: Date;
  Sites?: any;
}

export interface NetworkCard {
  Name: string;
  Mac: string;
}

