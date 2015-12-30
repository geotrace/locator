package locator

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// yandex описывает обработчик гео-запросов для сервиса Яндекс.
type yandex struct {
	apiKey string       // ключ для получения данных
	client *http.Client // HTTP-клиент
}

// Get передает данные на сервис Яндекс и возвращает разобранный ответ.
func (l *yandex) Get(req Request) (*Response, error) {
	// формируем данные для запроса в формате Yandex
	cells := make([]*yandexCell, len(req.CellTowers))
	for i, cell := range req.CellTowers {
		cells[i] = &yandexCell{
			MobileCountryCode: cell.MobileCountryCode,
			MobileNetworkCode: cell.MobileNetworkCode,
			CellId:            cell.CellId,
			LocationAreaCode:  cell.LocationAreaCode,
			SignalStrength:    cell.SignalStrength,
			Age:               cell.Age,
		}
	}
	wifis := make([]*yandexWiFi, len(req.WifiAccessPoints))
	for i, wifi := range req.WifiAccessPoints {
		wifis[i] = &yandexWiFi{
			MacAddress:     wifi.MacAddress,
			SignalStrength: wifi.SignalStrength,
			Age:            wifi.Age,
		}
	}
	request := yandexRequest{
		Common: yandexCommon{
			Version: "1.0",
			ApiKey:  l.apiKey,
		},
		GSMCells:     cells,
		WiFiNetworks: wifis,
	}
	if !IgnoreIPMethod && req.IPAddress != "" {
		request.IP = &yandexIP{AddressV4: req.IPAddress}
	}
	data, err := json.Marshal(&request)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(data))
	// для Яндекса это передается в виде запроса формы
	params := make(url.Values, 1)
	params.Set("json", string(data))
	resp, err := l.client.PostForm(Yandex, params)
	if err != nil {
		return nil, err
	}
	// согласно документации Яндекс не возвращает коды ошибок
	// так что это оставлено просто на всякий случай
	switch resp.StatusCode {
	case 200: // все хорошо — данные получены
	case 400: // неверный формат данных запроса или плохой ключ
		return nil, ErrBadRequest
	case 403: // исчерпан лимит запросов
		return nil, ErrForbidden
	case 404: // информация не найдена
		return nil, ErrNotFound
	default: // другая нехорошая ошибка
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}
	// декодируем ответ
	var yresp yandexRespose
	err = json.NewDecoder(resp.Body).Decode(&yresp)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	// проверяем текст ошибки и обрабатываем аналогично общему методу всех локаторов
	switch yresp.Error {
	case "": // все хорошо — ошибки нет — игнорируем
	case "JSON request is invalid":
		return nil, ErrBadRequest
	case "invalid api_key":
		return nil, ErrForbidden
	case "Location not found":
		return nil, ErrNotFound
	default: // другой текст ошибки
		return nil, errors.New(yresp.Error)
	}
	// проверяем, что метод определения не по IP
	if IgnoreIPMethod && yresp.Position.Type == "ip" {
		return nil, ErrNotFound
	}
	// формируем ответ
	response := Response{
		Location: Point{
			Lat: yresp.Position.Latitude,
			Lng: yresp.Position.Longitude,
		},
		Accuracy: yresp.Position.Precision,
	}
	return &response, nil
}

type yandexRequest struct {
	Common       yandexCommon  `json:"common"`
	GSMCells     []*yandexCell `json:"gsm_cells,omitempty"`
	WiFiNetworks []*yandexWiFi `json:"wifi_networks,omitempty"`
	IP           *yandexIP     `json:"ip,omitempty"`
}

type yandexIP struct {
	AddressV4 string `json:"address_v4"`
}

type yandexCell struct {
	MobileCountryCode uint16 `json:"countrycode"`
	MobileNetworkCode uint16 `json:"operatorid"`
	CellId            uint32 `json:"cellid"`
	LocationAreaCode  uint16 `json:"lac"`
	SignalStrength    int16  `json:"signal_strength,omitempty"`
	Age               uint32 `json:"age,omitempty"`
}
type yandexWiFi struct {
	MacAddress     string `json:"mac"`
	SignalStrength int16  `json:"signal_strength,omitempty"`
	Age            uint32 `json:"age,omitempty"`
}
type yandexCommon struct {
	Version string `json:"version"`
	ApiKey  string `json:"api_key"`
}

type yandexPosition struct {
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Altitude          float64 `json:"altitude"`
	Precision         float64 `json:"precision"`
	AltitudePrecision float64 `json:"altitude_precision"`
	Type              string  `json:"type"` // Способ определения местоположения: gsm, wifi, ip
}

type yandexRespose struct {
	Position *yandexPosition `json:"position"`
	Error    string          `json:"error"`
}
