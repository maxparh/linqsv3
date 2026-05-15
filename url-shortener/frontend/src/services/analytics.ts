import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

// Типы для ответов бэкенда
export interface AnalyticsOverview {
  total_clicks: number
  unique_clicks: number
  bounce_rate: number
  avg_time_on_site: number
}

export interface ClicksOverTime {
  labels: string[]
  values: number[]
}

export interface LocationStat {
  country: string
  country_code: string
  percent: number
}

export interface DeviceStat {
  name: string
  value: number
  percent: number
  color: string
}

export interface BounceRateStat {
  link_url: string
  bounce_rate: number
  total_sessions: number
}

export const analyticsAPI = {
  async getOverview(days: number = 7): Promise<AnalyticsOverview> {
    const response = await axios.get(`${API_URL}/analytics/overview`, {
      params: { days },
      headers: { 'Authorization': `Bearer ${localStorage.getItem('access_token')}` }
    })
    // Защита от null
    return response.data || { total_clicks: 0, unique_clicks: 0, bounce_rate: 0, avg_time_on_site: 0 }
  },

  async getClicksOverTime(days: number = 7): Promise<ClicksOverTime> {
    const response = await axios.get(`${API_URL}/analytics/clicks-over-time`, {
      params: { days },
      headers: { 'Authorization': `Bearer ${localStorage.getItem('access_token')}` }
    })
    return response.data || { labels: [], values: [] }
  },

  async getTopLocations(days: number = 7, limit: number = 5): Promise<LocationStat[]> {
    const response = await axios.get(`${API_URL}/analytics/top-locations`, {
      params: { days, limit },
      headers: { 'Authorization': `Bearer ${localStorage.getItem('access_token')}` }
    })
    return response.data || []
  },

  async getDeviceStats(days: number = 7): Promise<DeviceStat[]> {
    const response = await axios.get(`${API_URL}/analytics/devices`, {
      params: { days },
      headers: { 'Authorization': `Bearer ${localStorage.getItem('access_token')}` }
    })
    return response.data || []
  },

  async getBounceRate(days: number = 7): Promise<BounceRateStat[]> {
    const response = await axios.get(`${API_URL}/analytics/bounce-rate`, {
      params: { days },
      headers: { 'Authorization': `Bearer ${localStorage.getItem('access_token')}` }
    })
    return response.data || []
  }
}