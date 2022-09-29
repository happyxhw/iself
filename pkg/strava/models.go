package strava

import (
	"time"
)

// ActivityStats
// A set of rolled-up statistics and totals for an athlete
type ActivityStats struct {
	BiggestRideDistance       float64        `json:"biggest_ride_distance,omitempty" bson:"biggest_ride_distance"`               // The longest distance ridden by the athlete.
	BiggestClimbElevationGain float64        `json:"biggest_climb_elevation_gain,omitempty" bson:"biggest_climb_elevation_gain"` // The highest climb ridden by the athlete.
	RecentRideTotals          *ActivityTotal `json:"recent_ride_totals,omitempty" bson:"recent_ride_totals"`                     // The recent (last 4 weeks) ride stats for the athlete.
	RecentRunTotals           *ActivityTotal `json:"recent_run_totals,omitempty" bson:"recent_run_totals"`                       // The recent (last 4 weeks) run stats for the athlete.
	RecentSwimTotals          *ActivityTotal `json:"recent_swim_totals,omitempty" bson:"recent_swim_totals"`                     // The recent (last 4 weeks) swim stats for the athlete.
	YtdRideTotals             *ActivityTotal `json:"ytd_ride_totals,omitempty" bson:"ytd_ride_totals"`                           // The year to date ride stats for the athlete.
	YtdRunTotals              *ActivityTotal `json:"ytd_run_totals,omitempty" bson:"ytd_run_totals"`                             // The year to date run stats for the athlete.
	YtdSwimTotals             *ActivityTotal `json:"ytd_swim_totals,omitempty" bson:"ytd_swim_totals"`                           // The year to date swim stats for the athlete.
	AllRideTotals             *ActivityTotal `json:"all_ride_totals,omitempty" bson:"all_ride_totals"`                           // The all time ride stats for the athlete.
	AllRunTotals              *ActivityTotal `json:"all_run_totals,omitempty" bson:"all_run_totals"`                             // The all time run stats for the athlete.
	AllSwimTotals             *ActivityTotal `json:"all_swim_totals,omitempty" bson:"all_swim_totals"`                           // The all time swim stats for the athlete.
}

// ActivityTotal
// A roll-up of metrics pertaining to a set of activities. Values are in seconds and meters.
type ActivityTotal struct {
	Count            int     `json:"count,omitempty" bson:"count"`                         // The number of activities considered in this total.
	Distance         float64 `json:"distance,omitempty" bson:"distance"`                   // The total distance covered by the considered activities.
	MovingTime       int     `json:"moving_time,omitempty" bson:"moving_time"`             // The total moving time of the considered activities.
	ElapsedTime      int     `json:"elapsed_time,omitempty" bson:"elapsed_time"`           // The total elapsed time of the considered activities.
	ElevationGain    float64 `json:"elevation_gain,omitempty" bson:"elevation_gain"`       // The total elevation gain of the considered activities.
	AchievementCount int     `json:"achievement_count,omitempty" bson:"achievement_count"` // The total number of achievements of the considered activities.
}

// ActivityType
// An enumeration of the types an activity may have.
type ActivityType interface{}

// ActivityZone
type ActivityZone struct {
	Score               int                    `json:"score,omitempty" bson:"score"`                               // An instance of integer.
	DistributionBuckets *TimedZoneDistribution `json:"distribution_buckets,omitempty" bson:"distribution_buckets"` // An instance of #/TimedZoneDistribution.
	Type                string                 `json:"type,omitempty" bson:"type"`                                 // May take one of the following values: heartrate, power
	SensorBased         bool                   `json:"sensor_based,omitempty" bson:"sensor_based"`                 // An instance of boolean.
	Points              int                    `json:"points,omitempty" bson:"points"`                             // An instance of integer.
	CustomZones         bool                   `json:"custom_zones,omitempty" bson:"custom_zones"`                 // An instance of boolean.
	Max                 int                    `json:"max,omitempty" bson:"max"`                                   // An instance of integer.
}

// BaseStream
type BaseStream struct {
	OriginalSize int    `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
}

// Comment
type Comment struct {
	Id         int64           `json:"id,omitempty" bson:"id"`                   // The unique identifier of this comment
	ActivityId int64           `json:"activity_id,omitempty" bson:"activity_id"` // The identifier of the activity this comment is related to
	Text       string          `json:"text,omitempty" bson:"text"`               // The content of the comment
	Athlete    *SummaryAthlete `json:"athlete,omitempty" bson:"athlete"`         // An instance of SummaryAthlete.
	CreatedAt  time.Time       `json:"created_at,omitempty" bson:"created_at"`   // The time at which this comment was created.
}

// Error
type Error struct {
	Code     string `json:"code,omitempty" bson:"code"`         // The code associated with this error.
	Field    string `json:"field,omitempty" bson:"field"`       // The specific field or aspect of the resource associated with this error.
	Resource string `json:"resource,omitempty" bson:"resource"` // The type of resource associated with this error.
}

// ExplorerResponse
type ExplorerResponse struct {
	Segments *ExplorerSegment `json:"segments,omitempty" bson:"segments"` // The set of segments matching an explorer request
}

// ExplorerSegment
type ExplorerSegment struct {
	Id                int64   `json:"id,omitempty" bson:"id"`                                   // The unique identifier of this segment
	Name              string  `json:"name,omitempty" bson:"name"`                               // The name of this segment
	ClimbCategory     int     `json:"climb_category,omitempty" bson:"climb_category"`           // The category of the climb [0, 5]. Higher is harder ie. 5 is Hors catégorie, 0 is uncategorized in climb_category. If climb_category = 5, climb_category_desc = HC. If climb_category = 2, climb_category_desc = 3.
	ClimbCategoryDesc string  `json:"climb_category_desc,omitempty" bson:"climb_category_desc"` // The description for the category of the climb May take one of the following values: NC, 4, 3, 2, 1, HC
	AvgGrade          float64 `json:"avg_grade,omitempty" bson:"avg_grade"`                     // The segment's average grade, in percents
	StartLatlng       *LatLng `json:"start_latlng,omitempty" bson:"start_latlng"`               // An instance of LatLng.
	EndLatlng         *LatLng `json:"end_latlng,omitempty" bson:"end_latlng"`                   // An instance of LatLng.
	ElevDifference    float64 `json:"elev_difference,omitempty" bson:"elev_difference"`         // The segments's elevation difference, in meters
	Distance          float64 `json:"distance,omitempty" bson:"distance"`                       // The segment's distance, in meters
	Points            string  `json:"points,omitempty" bson:"points"`                           // The polyline of the segment
}

// Fault
// Encapsulates the errors that may be returned from the API.
type Fault struct {
	Errors  *Error `json:"errors,omitempty" bson:"errors"`   // The set of specific errors associated with this fault, if any.
	Message string `json:"message,omitempty" bson:"message"` // The message of the fault.
}

// HeartRateZoneRanges
type HeartRateZoneRanges struct {
	CustomZones bool        `json:"custom_zones,omitempty" bson:"custom_zones"` // Whether the athlete has set their own custom heart rate zones
	Zones       *ZoneRanges `json:"zones,omitempty" bson:"zones"`               // An instance of ZoneRanges.
}

// Lap
type Lap struct {
	Id                 int64         `json:"id,omitempty" bson:"id"`                                     // The unique identifier of this lap
	Activity           *MetaActivity `json:"activity,omitempty" bson:"activity"`                         // An instance of MetaActivity.
	Athlete            *MetaAthlete  `json:"athlete,omitempty" bson:"athlete"`                           // An instance of MetaAthlete.
	AverageCadence     float64       `json:"average_cadence,omitempty" bson:"average_cadence"`           // The lap's average cadence
	AverageSpeed       float64       `json:"average_speed,omitempty" bson:"average_speed"`               // The lap's average speed
	Distance           float64       `json:"distance,omitempty" bson:"distance"`                         // The lap's distance, in meters
	ElapsedTime        int           `json:"elapsed_time,omitempty" bson:"elapsed_time"`                 // The lap's elapsed time, in seconds
	StartIndex         int           `json:"start_index,omitempty" bson:"start_index"`                   // The start index of this effort in its activity's stream
	EndIndex           int           `json:"end_index,omitempty" bson:"end_index"`                       // The end index of this effort in its activity's stream
	LapIndex           int           `json:"lap_index,omitempty" bson:"lap_index"`                       // The index of this lap in the activity it belongs to
	MaxSpeed           float64       `json:"max_speed,omitempty" bson:"max_speed"`                       // The maximum speed of this lat, in meters per second
	MovingTime         int           `json:"moving_time,omitempty" bson:"moving_time"`                   // The lap's moving time, in seconds
	Name               string        `json:"name,omitempty" bson:"name"`                                 // The name of the lap
	PaceZone           int           `json:"pace_zone,omitempty" bson:"pace_zone"`                       // The athlete's pace zone during this lap
	Split              int           `json:"split,omitempty" bson:"split"`                               // An instance of integer.
	StartDate          time.Time     `json:"start_date,omitempty" bson:"start_date"`                     // The time at which the lap was started.
	StartDateLocal     time.Time     `json:"start_date_local,omitempty" bson:"start_date_local"`         // The time at which the lap was started in the local timezone.
	TotalElevationGain float64       `json:"total_elevation_gain,omitempty" bson:"total_elevation_gain"` // The elevation gain of this lap, in meters
}

// LatLng
// A collection of float objects. A pair of latitude/longitude coordinates, represented as an array of 2 floating point numbers.
type LatLng []float64

// MetaActivity
type MetaActivity struct {
	Id int64 `json:"id,omitempty" bson:"id"` // The unique identifier of the activity
}

// MetaAthlete
type MetaAthlete struct {
	ID int64 `json:"id,omitempty" bson:"id"` // The unique identifier of the athlete
}

// MetaClub
type MetaClub struct {
	Id            int64  `json:"id,omitempty" bson:"id"`                         // The club's unique identifier.
	ResourceState int    `json:"resource_state,omitempty" bson:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 -> "meta", 2 -> "summary", 3 -> "detail"
	Name          string `json:"name,omitempty" bson:"name"`                     // The club's name.
}

// PhotosSummary
type PhotosSummary struct {
	Count   int                    `json:"count,omitempty" bson:"count"`     // The number of photos
	Primary *PhotosSummary_primary `json:"primary,omitempty" bson:"primary"` // An instance of PhotosSummary_primary.
}

// PhotosSummary_primary
type PhotosSummary_primary struct {
	Id       int64  `json:"id,omitempty" bson:"id"`               // An instance of long.
	Source   int    `json:"source,omitempty" bson:"source"`       // An instance of integer.
	UniqueId string `json:"unique_id,omitempty" bson:"unique_id"` // An instance of string.
	Urls     string `json:"urls,omitempty" bson:"urls"`           // An instance of string.
}

// PolylineMap
type PolylineMap struct {
	Id              string `json:"id,omitempty" bson:"id"`                             // The identifier of the map
	Polyline        string `json:"polyline,omitempty" bson:"polyline"`                 // The polyline of the map, only returned on detailed representation of an object
	SummaryPolyline string `json:"summary_polyline,omitempty" bson:"summary_polyline"` // The summary polyline of the map
}

// PowerZoneRanges
type PowerZoneRanges struct {
	Zones *ZoneRanges `json:"zones,omitempty" bson:"zones"` // An instance of ZoneRanges.
}

// Route
type Route struct {
	Athlete             *SummaryAthlete `json:"athlete,omitempty" bson:"athlete"`                             // An instance of SummaryAthlete.
	Description         string          `json:"description,omitempty" bson:"description"`                     // The description of the route
	Distance            float64         `json:"distance,omitempty" bson:"distance"`                           // The route's distance, in meters
	ElevationGain       float64         `json:"elevation_gain,omitempty" bson:"elevation_gain"`               // The route's elevation gain.
	Id                  int64           `json:"id,omitempty" bson:"id"`                                       // The unique identifier of this route
	IdStr               string          `json:"id_str,omitempty" bson:"id_str"`                               // The unique identifier of the route in string format
	Map                 *PolylineMap    `json:"map,omitempty" bson:"map"`                                     // An instance of PolylineMap.
	Name                string          `json:"name,omitempty" bson:"name"`                                   // The name of this route
	Private             bool            `json:"private,omitempty" bson:"private"`                             // Whether this route is private
	Starred             bool            `json:"starred,omitempty" bson:"starred"`                             // Whether this route is starred by the logged-in athlete
	Timestamp           int             `json:"timestamp,omitempty" bson:"timestamp"`                         // An epoch timestamp of when the route was created
	Type                int             `json:"type,omitempty" bson:"type"`                                   // This route's type (1 for ride, 2 for runs)
	SubType             int             `json:"sub_type,omitempty" bson:"sub_type"`                           // This route's sub-type (1 for road, 2 for mountain bike, 3 for cross, 4 for trail, 5 for mixed)
	CreatedAt           time.Time       `json:"created_at,omitempty" bson:"created_at"`                       // The time at which the route was created
	UpdatedAt           time.Time       `json:"updated_at,omitempty" bson:"updated_at"`                       // The time at which the route was last updated
	EstimatedMovingTime int             `json:"estimated_moving_time,omitempty" bson:"estimated_moving_time"` // Estimated time in seconds for the authenticated athlete to complete route
	Segments            *SummarySegment `json:"segments,omitempty" bson:"segments"`                           // The segments traversed by this route
}

// RunningRace
type RunningRace struct {
	Id                    int64     `json:"id,omitempty" bson:"id"`                                         // The unique identifier of this race.
	Name                  string    `json:"name,omitempty" bson:"name"`                                     // The name of this race.
	RunningRaceType       int       `json:"running_race_type,omitempty" bson:"running_race_type"`           // The type of this race.
	Distance              float64   `json:"distance,omitempty" bson:"distance"`                             // The race's distance, in meters.
	StartDateLocal        time.Time `json:"start_date_local,omitempty" bson:"start_date_local"`             // The time at which the race begins started in the local timezone.
	City                  string    `json:"city,omitempty" bson:"city"`                                     // The name of the city in which the race is taking place.
	State                 string    `json:"state,omitempty" bson:"state"`                                   // The name of the state or geographical region in which the race is taking place.
	Country               string    `json:"country,omitempty" bson:"country"`                               // The name of the country in which the race is taking place.
	RouteIds              int64     `json:"route_ids,omitempty" bson:"route_ids"`                           // The set of routes that cover this race's course.
	MeasurementPreference string    `json:"measurement_preference,omitempty" bson:"measurement_preference"` // The unit system in which the race should be displayed. May take one of the following values: feet, meters
	Url                   string    `json:"url,omitempty" bson:"url"`                                       // The vanity URL of this race on Strava.
	WebsiteUrl            string    `json:"website_url,omitempty" bson:"website_url"`                       // The URL of this race's website.
}

// Split
type Split struct {
	AverageSpeed        float64 `json:"average_speed,omitempty" bson:"average_speed"`               // The average speed of this split, in meters per second
	Distance            float64 `json:"distance,omitempty" bson:"distance"`                         // The distance of this split, in meters
	ElapsedTime         int     `json:"elapsed_time,omitempty" bson:"elapsed_time"`                 // The elapsed time of this split, in seconds
	ElevationDifference float64 `json:"elevation_difference,omitempty" bson:"elevation_difference"` // The elevation difference of this split, in meters
	PaceZone            int     `json:"pace_zone,omitempty" bson:"pace_zone"`                       // The pacing zone of this split
	MovingTime          int     `json:"moving_time,omitempty" bson:"moving_time"`                   // The moving time of this split, in seconds
	Split               int     `json:"split,omitempty" bson:"split"`                               // N/A
}

// StreamSet
type StreamSet struct {
	Time           *TimeStream           `json:"time,omitempty" bson:"time"`                       // An instance of TimeStream.
	Distance       *DistanceStream       `json:"distance,omitempty" bson:"distance"`               // An instance of DistanceStream.
	Latlng         *LatLngStream         `json:"latlng,omitempty" bson:"latlng"`                   // An instance of LatLngStream.
	Altitude       *AltitudeStream       `json:"altitude,omitempty" bson:"altitude"`               // An instance of AltitudeStream.
	VelocitySmooth *SmoothVelocityStream `json:"velocity_smooth,omitempty" bson:"velocity_smooth"` // An instance of SmoothVelocityStream.
	Heartrate      *HeartrateStream      `json:"heartrate,omitempty" bson:"heartrate"`             // An instance of HeartrateStream.
	Cadence        *CadenceStream        `json:"cadence,omitempty" bson:"cadence"`                 // An instance of CadenceStream.
	Watts          *PowerStream          `json:"watts,omitempty" bson:"watts"`                     // An instance of PowerStream.
	Temp           *TemperatureStream    `json:"temp,omitempty" bson:"temp"`                       // An instance of TemperatureStream.
	Moving         *MovingStream         `json:"moving,omitempty" bson:"moving"`                   // An instance of MovingStream.
	GradeSmooth    *SmoothGradeStream    `json:"grade_smooth,omitempty" bson:"grade_smooth"`       // An instance of SmoothGradeStream.
}

// SummaryGear
type SummaryGear struct {
	Id            string  `json:"id,omitempty" bson:"id"`                         // The gear's unique identifier.
	ResourceState int     `json:"resource_state,omitempty" bson:"resource_state"` // Resource state, indicates level of detail. Possible values: 2 -> "summary", 3 -> "detail"
	Primary       bool    `json:"primary,omitempty" bson:"primary"`               // Whether this gear's is the owner's default one.
	Name          string  `json:"name,omitempty" bson:"name"`                     // The gear's name.
	Distance      float64 `json:"distance,omitempty" bson:"distance"`             // The distance logged with this gear.
}

// SummaryPRSegmentEffort
type SummaryPRSegmentEffort struct {
	PrActivityId  int64     `json:"pr_activity_id,omitempty" bson:"pr_activity_id"`   // The unique identifier of the activity related to the PR effort.
	PrElapsedTime int       `json:"pr_elapsed_time,omitempty" bson:"pr_elapsed_time"` // The elapsed time ot the PR effort.
	PrDate        time.Time `json:"pr_date,omitempty" bson:"pr_date"`                 // The time at which the PR effort was started.
	EffortCount   int       `json:"effort_count,omitempty" bson:"effort_count"`       // Number of efforts by the authenticated athlete on this segment.
}

// SummarySegment
type SummarySegment struct {
	Id                  int64                   `json:"id,omitempty" bson:"id"`                                       // The unique identifier of this segment
	Name                string                  `json:"name,omitempty" bson:"name"`                                   // The name of this segment
	ActivityType        string                  `json:"activity_type,omitempty" bson:"activity_type"`                 // May take one of the following values: Ride, Run
	Distance            float64                 `json:"distance,omitempty" bson:"distance"`                           // The segment's distance, in meters
	AverageGrade        float64                 `json:"average_grade,omitempty" bson:"average_grade"`                 // The segment's average grade, in percents
	MaximumGrade        float64                 `json:"maximum_grade,omitempty" bson:"maximum_grade"`                 // The segments's maximum grade, in percents
	ElevationHigh       float64                 `json:"elevation_high,omitempty" bson:"elevation_high"`               // The segments's highest elevation, in meters
	ElevationLow        float64                 `json:"elevation_low,omitempty" bson:"elevation_low"`                 // The segments's lowest elevation, in meters
	StartLatlng         *LatLng                 `json:"start_latlng,omitempty" bson:"start_latlng"`                   // An instance of LatLng.
	EndLatlng           *LatLng                 `json:"end_latlng,omitempty" bson:"end_latlng"`                       // An instance of LatLng.
	ClimbCategory       int                     `json:"climb_category,omitempty" bson:"climb_category"`               // The category of the climb [0, 5]. Higher is harder ie. 5 is Hors catégorie, 0 is uncategorized in climb_category.
	City                string                  `json:"city,omitempty" bson:"city"`                                   // The segments's city.
	State               string                  `json:"state,omitempty" bson:"state"`                                 // The segments's state or geographical region.
	Country             string                  `json:"country,omitempty" bson:"country"`                             // The segment's country.
	Private             bool                    `json:"private,omitempty" bson:"private"`                             // Whether this segment is private.
	AthletePrEffort     *SummarySegmentEffort   `json:"athlete_pr_effort,omitempty" bson:"athlete_pr_effort"`         // An instance of SummarySegmentEffort.
	AthleteSegmentStats *SummaryPRSegmentEffort `json:"athlete_segment_stats,omitempty" bson:"athlete_segment_stats"` // An instance of SummaryPRSegmentEffort.
}

// SummarySegmentEffort
type SummarySegmentEffort struct {
	Id             int64     `json:"id,omitempty" bson:"id"`                             // The unique identifier of this effort
	ActivityId     int64     `json:"activity_id,omitempty" bson:"activity_id"`           // The unique identifier of the activity related to this effort
	ElapsedTime    int       `json:"elapsed_time,omitempty" bson:"elapsed_time"`         // The effort's elapsed time
	StartDate      time.Time `json:"start_date,omitempty" bson:"start_date"`             // The time at which the effort was started.
	StartDateLocal time.Time `json:"start_date_local,omitempty" bson:"start_date_local"` // The time at which the effort was started in the local timezone.
	Distance       float64   `json:"distance,omitempty" bson:"distance"`                 // The effort's distance in meters
	IsKom          bool      `json:"is_kom,omitempty" bson:"is_kom"`                     // Whether this effort is the current best on the leaderboard
}

// TimedZoneDistribution
// A collection of #/TimedZoneRange objects. Stores the exclusive ranges representing zones and the time spent in each.
type TimedZoneDistribution []*TimedZoneRange

// UpdatableActivity
type UpdatableActivity struct {
	Commute     bool          `json:"commute,omitempty" bson:"commute"`         // Whether this activity is a commute
	Trainer     bool          `json:"trainer,omitempty" bson:"trainer"`         // Whether this activity was recorded on a training machine
	Description string        `json:"description,omitempty" bson:"description"` // The description of the activity
	Name        string        `json:"name,omitempty" bson:"name"`               // The name of the activity
	Type        *ActivityType `json:"type,omitempty" bson:"type"`               // An instance of ActivityType.
	GearId      string        `json:"gear_id,omitempty" bson:"gear_id"`         // Identifier for the gear associated with the activity. ‘none’ clears gear from activity
}

// Upload
type Upload struct {
	Id         int64  `json:"id,omitempty" bson:"id"`                   // The unique identifier of the upload
	IdStr      string `json:"id_str,omitempty" bson:"id_str"`           // The unique identifier of the upload in string format
	ExternalId string `json:"external_id,omitempty" bson:"external_id"` // The external identifier of the upload
	Error      string `json:"error,omitempty" bson:"error"`             // The error associated with this upload
	Status     string `json:"status,omitempty" bson:"status"`           // The status of this upload
	ActivityId int64  `json:"activity_id,omitempty" bson:"activity_id"` // The identifier of the activity this upload resulted into
}

// ZoneRange
type ZoneRange struct {
	Min int `json:"min,omitempty" bson:"min"` // The minimum value in the range.
	Max int `json:"max,omitempty" bson:"max"` // The maximum value in the range.
}

// ZoneRanges
// A collection of ZoneRange objects.
type ZoneRanges []*ZoneRange

// Zones
type Zones struct {
	HeartRate *HeartRateZoneRanges `json:"heart_rate,omitempty" bson:"heart_rate"` // An instance of HeartRateZoneRanges.
	Power     *PowerZoneRanges     `json:"power,omitempty" bson:"power"`           // An instance of PowerZoneRanges.
}

// AltitudeStream
type AltitudeStream struct {
	OriginalSize int       `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string    `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string    `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []float64 `json:"data,omitempty" bson:"data"`                   // The sequence of altitude values for this stream, in meters
}

// CadenceStream
type CadenceStream struct {
	OriginalSize int    `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []int  `json:"data,omitempty" bson:"data"`                   // The sequence of cadence values for this stream, in rotations per minute
}

// DetailedGear
type DetailedGear struct {
	Id            string  `json:"id,omitempty" bson:"id"`                         // The gear's unique identifier.
	ResourceState int     `json:"resource_state,omitempty" bson:"resource_state"` // Resource state, indicates level of detail. Possible values: 2 -> "summary", 3 -> "detail"
	Primary       bool    `json:"primary,omitempty" bson:"primary"`               // Whether this gear's is the owner's default one.
	Name          string  `json:"name,omitempty" bson:"name"`                     // The gear's name.
	Distance      float64 `json:"distance,omitempty" bson:"distance"`             // The distance logged with this gear.
	BrandName     string  `json:"brand_name,omitempty" bson:"brand_name"`         // The gear's brand name.
	ModelName     string  `json:"model_name,omitempty" bson:"model_name"`         // The gear's model name.
	FrameType     int     `json:"frame_type,omitempty" bson:"frame_type"`         // The gear's frame type (bike only).
	Description   string  `json:"description,omitempty" bson:"description"`       // The gear's description.
}

// DetailedSegment
type DetailedSegment struct {
	Id                  int64                   `json:"id,omitempty" bson:"id"`                                       // The unique identifier of this segment
	Name                string                  `json:"name,omitempty" bson:"name"`                                   // The name of this segment
	ActivityType        string                  `json:"activity_type,omitempty" bson:"activity_type"`                 // May take one of the following values: Ride, Run
	Distance            float64                 `json:"distance,omitempty" bson:"distance"`                           // The segment's distance, in meters
	AverageGrade        float64                 `json:"average_grade,omitempty" bson:"average_grade"`                 // The segment's average grade, in percents
	MaximumGrade        float64                 `json:"maximum_grade,omitempty" bson:"maximum_grade"`                 // The segments's maximum grade, in percents
	ElevationHigh       float64                 `json:"elevation_high,omitempty" bson:"elevation_high"`               // The segments's highest elevation, in meters
	ElevationLow        float64                 `json:"elevation_low,omitempty" bson:"elevation_low"`                 // The segments's lowest elevation, in meters
	StartLatlng         *LatLng                 `json:"start_latlng,omitempty" bson:"start_latlng"`                   // An instance of LatLng.
	EndLatlng           *LatLng                 `json:"end_latlng,omitempty" bson:"end_latlng"`                       // An instance of LatLng.
	ClimbCategory       int                     `json:"climb_category,omitempty" bson:"climb_category"`               // The category of the climb [0, 5]. Higher is harder ie. 5 is Hors catégorie, 0 is uncategorized in climb_category.
	City                string                  `json:"city,omitempty" bson:"city"`                                   // The segments's city.
	State               string                  `json:"state,omitempty" bson:"state"`                                 // The segments's state or geographical region.
	Country             string                  `json:"country,omitempty" bson:"country"`                             // The segment's country.
	Private             bool                    `json:"private,omitempty" bson:"private"`                             // Whether this segment is private.
	AthletePrEffort     *SummarySegmentEffort   `json:"athlete_pr_effort,omitempty" bson:"athlete_pr_effort"`         // An instance of SummarySegmentEffort.
	AthleteSegmentStats *SummaryPRSegmentEffort `json:"athlete_segment_stats,omitempty" bson:"athlete_segment_stats"` // An instance of SummaryPRSegmentEffort.
	CreatedAt           time.Time               `json:"created_at,omitempty" bson:"created_at"`                       // The time at which the segment was created.
	UpdatedAt           time.Time               `json:"updated_at,omitempty" bson:"updated_at"`                       // The time at which the segment was last updated.
	TotalElevationGain  float64                 `json:"total_elevation_gain,omitempty" bson:"total_elevation_gain"`   // The segment's total elevation gain.
	Map                 *PolylineMap            `json:"map,omitempty" bson:"map"`                                     // An instance of PolylineMap.
	EffortCount         int                     `json:"effort_count,omitempty" bson:"effort_count"`                   // The total number of efforts for this segment
	AthleteCount        int                     `json:"athlete_count,omitempty" bson:"athlete_count"`                 // The number of unique athletes who have an effort for this segment
	Hazardous           bool                    `json:"hazardous,omitempty" bson:"hazardous"`                         // Whether this segment is considered hazardous
	StarCount           int                     `json:"star_count,omitempty" bson:"star_count"`                       // The number of stars for this segment
}

// DetailedSegmentEffort
type DetailedSegmentEffort struct {
	Id               int64           `json:"id,omitempty" bson:"id"`                               // The unique identifier of this effort
	ActivityId       int64           `json:"activity_id,omitempty" bson:"activity_id"`             // The unique identifier of the activity related to this effort
	ElapsedTime      int             `json:"elapsed_time,omitempty" bson:"elapsed_time"`           // The effort's elapsed time
	StartDate        time.Time       `json:"start_date,omitempty" bson:"start_date"`               // The time at which the effort was started.
	StartDateLocal   time.Time       `json:"start_date_local,omitempty" bson:"start_date_local"`   // The time at which the effort was started in the local timezone.
	Distance         float64         `json:"distance,omitempty" bson:"distance"`                   // The effort's distance in meters
	IsKom            bool            `json:"is_kom,omitempty" bson:"is_kom"`                       // Whether this effort is the current best on the leaderboard
	Name             string          `json:"name,omitempty" bson:"name"`                           // The name of the segment on which this effort was performed
	Activity         *MetaActivity   `json:"activity,omitempty" bson:"activity"`                   // An instance of MetaActivity.
	Athlete          *MetaAthlete    `json:"athlete,omitempty" bson:"athlete"`                     // An instance of MetaAthlete.
	MovingTime       int             `json:"moving_time,omitempty" bson:"moving_time"`             // The effort's moving time
	StartIndex       int             `json:"start_index,omitempty" bson:"start_index"`             // The start index of this effort in its activity's stream
	EndIndex         int             `json:"end_index,omitempty" bson:"end_index"`                 // The end index of this effort in its activity's stream
	AverageCadence   float64         `json:"average_cadence,omitempty" bson:"average_cadence"`     // The effort's average cadence
	AverageWatts     float64         `json:"average_watts,omitempty" bson:"average_watts"`         // The average wattage of this effort
	DeviceWatts      bool            `json:"device_watts,omitempty" bson:"device_watts"`           // For riding efforts, whether the wattage was reported by a dedicated recording device
	AverageHeartrate float64         `json:"average_heartrate,omitempty" bson:"average_heartrate"` // The heart heart rate of the athlete during this effort
	MaxHeartrate     float64         `json:"max_heartrate,omitempty" bson:"max_heartrate"`         // The maximum heart rate of the athlete during this effort
	Segment          *SummarySegment `json:"segment,omitempty" bson:"segment"`                     // An instance of SummarySegment.
	KomRank          int             `json:"kom_rank,omitempty" bson:"kom_rank"`                   // The rank of the effort on the global leaderboard if it belongs in the top 10 at the time of upload
	PrRank           int             `json:"pr_rank,omitempty" bson:"pr_rank"`                     // The rank of the effort on the athlete's leaderboard if it belongs in the top 3 at the time of upload
	Hidden           bool            `json:"hidden,omitempty" bson:"hidden"`                       // Whether this effort should be hidden when viewed within an activity
}

// DistanceStream
type DistanceStream struct {
	OriginalSize int       `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string    `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string    `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []float64 `json:"data,omitempty" bson:"data"`                   // The sequence of distance values for this stream, in meters
}

// HeartrateStream
type HeartrateStream struct {
	OriginalSize int    `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []int  `json:"data,omitempty" bson:"data"`                   // The sequence of heart rate values for this stream, in beats per minute
}

// LatLngStream
type LatLngStream struct {
	OriginalSize int       `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string    `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string    `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []*LatLng `json:"data,omitempty" bson:"data"`                   // The sequence of lat/long values for this stream
}

// MovingStream
type MovingStream struct {
	OriginalSize int    `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []bool `json:"data,omitempty" bson:"data"`                   // The sequence of moving values for this stream, as boolean values
}

// PowerStream
type PowerStream struct {
	OriginalSize int    `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []int  `json:"data,omitempty" bson:"data"`                   // The sequence of power values for this stream, in watts
}

// SmoothGradeStream
type SmoothGradeStream struct {
	OriginalSize int       `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string    `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string    `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []float64 `json:"data,omitempty" bson:"data"`                   // The sequence of grade values for this stream, as percents of a grade
}

// SmoothVelocityStream
type SmoothVelocityStream struct {
	OriginalSize int       `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string    `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string    `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []float64 `json:"data,omitempty" bson:"data"`                   // The sequence of velocity values for this stream, in meters per second
}

// SummaryActivity
type SummaryActivity struct {
	Id                   int64         `json:"id,omitempty" bson:"id"`                                         // The unique identifier of the activity
	ExternalId           string        `json:"external_id,omitempty" bson:"external_id"`                       // The identifier provided at upload time
	UploadId             int64         `json:"upload_id,omitempty" bson:"upload_id"`                           // The identifier of the upload that resulted in this activity
	Athlete              *MetaAthlete  `json:"athlete,omitempty" bson:"athlete"`                               // An instance of MetaAthlete.
	Name                 string        `json:"name,omitempty" bson:"name"`                                     // The name of the activity
	Distance             float64       `json:"distance,omitempty" bson:"distance"`                             // The activity's distance, in meters
	MovingTime           int           `json:"moving_time,omitempty" bson:"moving_time"`                       // The activity's moving time, in seconds
	ElapsedTime          int           `json:"elapsed_time,omitempty" bson:"elapsed_time"`                     // The activity's elapsed time, in seconds
	TotalElevationGain   float64       `json:"total_elevation_gain,omitempty" bson:"total_elevation_gain"`     // The activity's total elevation gain.
	ElevHigh             float64       `json:"elev_high,omitempty" bson:"elev_high"`                           // The activity's highest elevation, in meters
	ElevLow              float64       `json:"elev_low,omitempty" bson:"elev_low"`                             // The activity's lowest elevation, in meters
	Type                 *ActivityType `json:"type,omitempty" bson:"type"`                                     // An instance of ActivityType.
	StartDate            time.Time     `json:"start_date,omitempty" bson:"start_date"`                         // The time at which the activity was started.
	StartDateLocal       time.Time     `json:"start_date_local,omitempty" bson:"start_date_local"`             // The time at which the activity was started in the local timezone.
	Timezone             string        `json:"timezone,omitempty" bson:"timezone"`                             // The timezone of the activity
	StartLatlng          *LatLng       `json:"start_latlng,omitempty" bson:"start_latlng"`                     // An instance of LatLng.
	EndLatlng            *LatLng       `json:"end_latlng,omitempty" bson:"end_latlng"`                         // An instance of LatLng.
	AchievementCount     int           `json:"achievement_count,omitempty" bson:"achievement_count"`           // The number of achievements gained during this activity
	KudosCount           int           `json:"kudos_count,omitempty" bson:"kudos_count"`                       // The number of kudos given for this activity
	CommentCount         int           `json:"comment_count,omitempty" bson:"comment_count"`                   // The number of comments for this activity
	AthleteCount         int           `json:"athlete_count,omitempty" bson:"athlete_count"`                   // The number of athletes for taking part in a group activity
	PhotoCount           int           `json:"photo_count,omitempty" bson:"photo_count"`                       // The number of Instagram photos for this activity
	TotalPhotoCount      int           `json:"total_photo_count,omitempty" bson:"total_photo_count"`           // The number of Instagram and Strava photos for this activity
	Map                  *PolylineMap  `json:"map,omitempty" bson:"map"`                                       // An instance of PolylineMap.
	Trainer              bool          `json:"trainer,omitempty" bson:"trainer"`                               // Whether this activity was recorded on a training machine
	Commute              bool          `json:"commute,omitempty" bson:"commute"`                               // Whether this activity is a commute
	Manual               bool          `json:"manual,omitempty" bson:"manual"`                                 // Whether this activity was created manually
	Private              bool          `json:"private,omitempty" bson:"private"`                               // Whether this activity is private
	Flagged              bool          `json:"flagged,omitempty" bson:"flagged"`                               // Whether this activity is flagged
	WorkoutType          int           `json:"workout_type,omitempty" bson:"workout_type"`                     // The activity's workout type
	UploadIdStr          string        `json:"upload_id_str,omitempty" bson:"upload_id_str"`                   // The unique identifier of the upload in string format
	AverageSpeed         float64       `json:"average_speed,omitempty" bson:"average_speed"`                   // The activity's average speed, in meters per second
	MaxSpeed             float64       `json:"max_speed,omitempty" bson:"max_speed"`                           // The activity's max speed, in meters per second
	HasKudoed            bool          `json:"has_kudoed,omitempty" bson:"has_kudoed"`                         // Whether the logged-in athlete has kudoed this activity
	GearId               string        `json:"gear_id,omitempty" bson:"gear_id"`                               // The id of the gear for the activity
	Kilojoules           float64       `json:"kilojoules,omitempty" bson:"kilojoules"`                         // The total work done in kilojoules during this activity. Rides only
	AverageWatts         float64       `json:"average_watts,omitempty" bson:"average_watts"`                   // Average power output in watts during this activity. Rides only
	DeviceWatts          bool          `json:"device_watts,omitempty" bson:"device_watts"`                     // Whether the watts are from a power meter, false if estimated
	MaxWatts             int           `json:"max_watts,omitempty" bson:"max_watts"`                           // Rides with power meter data only
	WeightedAverageWatts int           `json:"weighted_average_watts,omitempty" bson:"weighted_average_watts"` // Similar to Normalized Power. Rides with power meter data only
}

// SummaryAthlete
type SummaryAthlete struct {
	Id            int64     `json:"id,omitempty" bson:"id"`                         // The unique identifier of the athlete
	Username      string    `json:"username" bson:"username"`                       // user name
	ResourceState int       `json:"resource_state,omitempty" bson:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 -> "meta", 2 -> "summary", 3 -> "detail"
	Firstname     string    `json:"firstname,omitempty" bson:"firstname"`           // The athlete's first name.
	Lastname      string    `json:"lastname,omitempty" bson:"lastname"`             // The athlete's last name.
	ProfileMedium string    `json:"profile_medium,omitempty" bson:"profile_medium"` // URL to a 62x62 pixel profile picture.
	Profile       string    `json:"profile,omitempty" bson:"profile"`               // URL to a 124x124 pixel profile picture.
	City          string    `json:"city,omitempty" bson:"city"`                     // The athlete's city.
	State         string    `json:"state,omitempty" bson:"state"`                   // The athlete's state or geographical region.
	Country       string    `json:"country,omitempty" bson:"country"`               // The athlete's country.
	Sex           string    `json:"sex,omitempty" bson:"sex"`                       // The athlete's sex. May take one of the following values: M, F
	Premium       bool      `json:"premium,omitempty" bson:"premium"`               // Deprecated.  Use summit field instead. Whether the athlete has any Summit subscription.
	Summit        bool      `json:"summit,omitempty" bson:"summit"`                 // Whether the athlete has any Summit subscription.
	CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at"`         // The time at which the athlete was created.
	UpdatedAt     time.Time `json:"updated_at,omitempty" bson:"updated_at"`         // The time at which the athlete was last updated.
}

// SummaryClub
type SummaryClub struct {
	Id              int64  `json:"id,omitempty" bson:"id"`                               // The club's unique identifier.
	ResourceState   int    `json:"resource_state,omitempty" bson:"resource_state"`       // Resource state, indicates level of detail. Possible values: 1 -> "meta", 2 -> "summary", 3 -> "detail"
	Name            string `json:"name,omitempty" bson:"name"`                           // The club's name.
	ProfileMedium   string `json:"profile_medium,omitempty" bson:"profile_medium"`       // URL to a 60x60 pixel profile picture.
	CoverPhoto      string `json:"cover_photo,omitempty" bson:"cover_photo"`             // URL to a ~1185x580 pixel cover photo.
	CoverPhotoSmall string `json:"cover_photo_small,omitempty" bson:"cover_photo_small"` // URL to a ~360x176  pixel cover photo.
	SportType       string `json:"sport_type,omitempty" bson:"sport_type"`               // May take one of the following values: cycling, running, triathlon, other
	City            string `json:"city,omitempty" bson:"city"`                           // The club's city.
	State           string `json:"state,omitempty" bson:"state"`                         // The club's state or geographical region.
	Country         string `json:"country,omitempty" bson:"country"`                     // The club's country.
	Private         bool   `json:"private,omitempty" bson:"private"`                     // Whether the club is private.
	MemberCount     int    `json:"member_count,omitempty" bson:"member_count"`           // The club's member count.
	Featured        bool   `json:"featured,omitempty" bson:"featured"`                   // Whether the club is featured or not.
	Verified        bool   `json:"verified,omitempty" bson:"verified"`                   // Whether the club is verified or not.
	Url             string `json:"url,omitempty" bson:"url"`                             // The club's vanity URL.
}

// TemperatureStream
type TemperatureStream struct {
	OriginalSize int    `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []int  `json:"data,omitempty" bson:"data"`                   // The sequence of temperature values for this stream, in celsius degrees
}

// TimeStream
type TimeStream struct {
	OriginalSize int    `json:"original_size,omitempty" bson:"original_size"` // The number of data points in this stream
	Resolution   string `json:"resolution,omitempty" bson:"resolution"`       // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string `json:"series_type,omitempty" bson:"series_type"`     // The base series used in the case the stream was downsampled May take one of the following values: distance, time
	Data         []int  `json:"data,omitempty" bson:"data"`                   // The sequence of time values for this stream, in seconds
}

// TimedZoneRange
// A union type representing the time spent in a given zone.
type TimedZoneRange struct {
	Min  int `json:"min,omitempty" bson:"min"`   // The minimum value in the range.
	Max  int `json:"max,omitempty" bson:"max"`   // The maximum value in the range.
	Time int `json:"time,omitempty" bson:"time"` // The number of seconds spent in this zone
}

// DetailedActivity
type DetailedActivity struct {
	ID                   int64                    `json:"id,omitempty" bson:"id"`                                         // The unique identifier of the activity
	ExternalID           string                   `json:"external_id,omitempty" bson:"external_id"`                       // The identifier provided at upload time
	UploadID             int64                    `json:"upload_id,omitempty" bson:"upload_id"`                           // The identifier of the upload that resulted in this activity
	Athlete              *MetaAthlete             `json:"athlete,omitempty" bson:"athlete"`                               // An instance of MetaAthlete.
	Name                 string                   `json:"name,omitempty" bson:"name"`                                     // The name of the activity
	Distance             float64                  `json:"distance,omitempty" bson:"distance"`                             // The activity's distance, in meters
	MovingTime           int                      `json:"moving_time,omitempty" bson:"moving_time"`                       // The activity's moving time, in seconds
	ElapsedTime          int                      `json:"elapsed_time,omitempty" bson:"elapsed_time"`                     // The activity's elapsed time, in seconds
	TotalElevationGain   float64                  `json:"total_elevation_gain,omitempty" bson:"total_elevation_gain"`     // The activity's total elevation gain.
	ElevHigh             float64                  `json:"elev_high,omitempty" bson:"elev_high"`                           // The activity's highest elevation, in meters
	ElevLow              float64                  `json:"elev_low,omitempty" bson:"elev_low"`                             // The activity's lowest elevation, in meters
	Type                 string                   `json:"type,omitempty" bson:"type"`                                     // An instance of ActivityType.
	StartDate            time.Time                `json:"start_date,omitempty" bson:"start_date"`                         // The time at which the activity was started.
	StartDateLocal       time.Time                `json:"start_date_local,omitempty" bson:"start_date_local"`             // The time at which the activity was started in the local timezone.
	Timezone             string                   `json:"timezone,omitempty" bson:"timezone"`                             // The timezone of the activity
	StartLatlng          *LatLng                  `json:"start_latlng,omitempty" bson:"start_latlng"`                     // An instance of LatLng.
	EndLatlng            *LatLng                  `json:"end_latlng,omitempty" bson:"end_latlng"`                         // An instance of LatLng.
	AchievementCount     int                      `json:"achievement_count,omitempty" bson:"achievement_count"`           // The number of achievements gained during this activity
	KudosCount           int                      `json:"kudos_count,omitempty" bson:"kudos_count"`                       // The number of kudos given for this activity
	CommentCount         int                      `json:"comment_count,omitempty" bson:"comment_count"`                   // The number of comments for this activity
	AthleteCount         int                      `json:"athlete_count,omitempty" bson:"athlete_count"`                   // The number of athletes for taking part in a group activity
	PhotoCount           int                      `json:"photo_count,omitempty" bson:"photo_count"`                       // The number of Instagram photos for this activity
	TotalPhotoCount      int                      `json:"total_photo_count,omitempty" bson:"total_photo_count"`           // The number of Instagram and Strava photos for this activity
	Map                  *PolylineMap             `json:"map,omitempty" bson:"map"`                                       // An instance of PolylineMap.
	Trainer              bool                     `json:"trainer,omitempty" bson:"trainer"`                               // Whether this activity was recorded on a training machine
	Commute              bool                     `json:"commute,omitempty" bson:"commute"`                               // Whether this activity is a commute
	Manual               bool                     `json:"manual,omitempty" bson:"manual"`                                 // Whether this activity was created manually
	Private              bool                     `json:"private,omitempty" bson:"private"`                               // Whether this activity is private
	Flagged              bool                     `json:"flagged,omitempty" bson:"flagged"`                               // Whether this activity is flagged
	WorkoutType          int                      `json:"workout_type,omitempty" bson:"workout_type"`                     // The activity's workout type
	UploadIdStr          string                   `json:"upload_id_str,omitempty" bson:"upload_id_str"`                   // The unique identifier of the upload in string format
	AverageSpeed         float64                  `json:"average_speed,omitempty" bson:"average_speed"`                   // The activity's average speed, in meters per second
	MaxSpeed             float64                  `json:"max_speed,omitempty" bson:"max_speed"`                           // The activity's max speed, in meters per second
	HasKudoed            bool                     `json:"has_kudoed,omitempty" bson:"has_kudoed"`                         // Whether the logged-in athlete has kudoed this activity
	GearId               string                   `json:"gear_id,omitempty" bson:"gear_id"`                               // The id of the gear for the activity
	Kilojoules           float64                  `json:"kilojoules,omitempty" bson:"kilojoules"`                         // The total work done in kilojoules during this activity. Rides only
	AverageWatts         float64                  `json:"average_watts,omitempty" bson:"average_watts"`                   // Average power output in watts during this activity. Rides only
	DeviceWatts          bool                     `json:"device_watts,omitempty" bson:"device_watts"`                     // Whether the watts are from a power meter, false if estimated
	MaxWatts             int                      `json:"max_watts,omitempty" bson:"max_watts"`                           // Rides with power meter data only
	WeightedAverageWatts int                      `json:"weighted_average_watts,omitempty" bson:"weighted_average_watts"` // Similar to Normalized Power. Rides with power meter data only
	Description          string                   `json:"description,omitempty" bson:"description"`                       // The description of the activity
	Photos               *PhotosSummary           `json:"photos,omitempty" bson:"photos"`                                 // An instance of PhotosSummary.
	Gear                 *SummaryGear             `json:"gear,omitempty" bson:"gear"`                                     // An instance of SummaryGear.
	Calories             float64                  `json:"calories,omitempty" bson:"calories"`                             // The number of kilocalories consumed during this activity
	SegmentEfforts       []*DetailedSegmentEffort `json:"segment_efforts,omitempty" bson:"segment_efforts"`               // A collection of DetailedSegmentEffort objects.
	DeviceName           string                   `json:"device_name,omitempty" bson:"device_name"`                       // The name of the device used to record the activity
	EmbedToken           string                   `json:"embed_token,omitempty" bson:"embed_token"`                       // The token used to embed a Strava activity
	SplitsMetric         []*Split                 `json:"splits_metric,omitempty" bson:"splits_metric"`                   // The splits of this activity in metric units (for runs)
	SplitsStandard       []*Split                 `json:"splits_standard,omitempty" bson:"splits_standard"`               // The splits of this activity in imperial units (for runs)
	Laps                 []*Lap                   `json:"laps,omitempty" bson:"laps"`                                     // A collection of Lap objects.
	BestEfforts          []*DetailedSegmentEffort `json:"best_efforts,omitempty" bson:"best_efforts"`
	AverageHeartrate     float64                  `json:"average_heartrate,omitempty" bson:"average_heartrate"` // The heart heart rate of the athlete during this effort
	MaxHeartrate         float64                  `json:"max_heartrate,omitempty" bson:"max_heartrate"`         // The maximum heart rate of the athlete during this effort
}

// DetailedAthlete
type DetailedAthlete struct {
	Id                    int64        `json:"id,omitempty" bson:"id"`                                         // The unique identifier of the athlete
	ResourceState         int          `json:"resource_state,omitempty" bson:"resource_state"`                 // Resource state, indicates level of detail. Possible values: 1 -> "meta", 2 -> "summary", 3 -> "detail"
	Firstname             string       `json:"firstname,omitempty" bson:"firstname"`                           // The athlete's first name.
	Lastname              string       `json:"lastname,omitempty" bson:"lastname"`                             // The athlete's last name.
	ProfileMedium         string       `json:"profile_medium,omitempty" bson:"profile_medium"`                 // URL to a 62x62 pixel profile picture.
	Profile               string       `json:"profile,omitempty" bson:"profile"`                               // URL to a 124x124 pixel profile picture.
	City                  string       `json:"city,omitempty" bson:"city"`                                     // The athlete's city.
	State                 string       `json:"state,omitempty" bson:"state"`                                   // The athlete's state or geographical region.
	Country               string       `json:"country,omitempty" bson:"country"`                               // The athlete's country.
	Sex                   string       `json:"sex,omitempty" bson:"sex"`                                       // The athlete's sex. May take one of the following values: M, F
	Premium               bool         `json:"premium,omitempty" bson:"premium"`                               // Deprecated.  Use summit field instead. Whether the athlete has any Summit subscription.
	Summit                bool         `json:"summit,omitempty" bson:"summit"`                                 // Whether the athlete has any Summit subscription.
	CreatedAt             time.Time    `json:"created_at,omitempty" bson:"created_at"`                         // The time at which the athlete was created.
	UpdatedAt             time.Time    `json:"updated_at,omitempty" bson:"updated_at"`                         // The time at which the athlete was last updated.
	FollowerCount         int          `json:"follower_count,omitempty" bson:"follower_count"`                 // The athlete's follower count.
	FriendCount           int          `json:"friend_count,omitempty" bson:"friend_count"`                     // The athlete's friend count.
	MeasurementPreference string       `json:"measurement_preference,omitempty" bson:"measurement_preference"` // The athlete's preferred unit systex. May take one of the following values: feet, meters
	Ftp                   int          `json:"ftp,omitempty" bson:"ftp"`                                       // The athlete's FTP (Functional Threshold Power).
	Weight                float64      `json:"weight,omitempty" bson:"weight"`                                 // The athlete's weight.
	Clubs                 *SummaryClub `json:"clubs,omitempty" bson:"clubs"`                                   // The athlete's clubs.
	Bikes                 *SummaryGear `json:"bikes,omitempty" bson:"bikes"`                                   // The athlete's bikes.
	Shoes                 *SummaryGear `json:"shoes,omitempty" bson:"shoes"`                                   // The athlete's shoes.
}

// DetailedClub
type DetailedClub struct {
	Id              int64  `json:"id,omitempty" bson:"id"`                               // The club's unique identifier.
	ResourceState   int    `json:"resource_state,omitempty" bson:"resource_state"`       // Resource state, indicates level of detail. Possible values: 1 -> "meta", 2 -> "summary", 3 -> "detail"
	Name            string `json:"name,omitempty" bson:"name"`                           // The club's name.
	ProfileMedium   string `json:"profile_medium,omitempty" bson:"profile_medium"`       // URL to a 60x60 pixel profile picture.
	CoverPhoto      string `json:"cover_photo,omitempty" bson:"cover_photo"`             // URL to a ~1185x580 pixel cover photo.
	CoverPhotoSmall string `json:"cover_photo_small,omitempty" bson:"cover_photo_small"` // URL to a ~360x176  pixel cover photo.
	SportType       string `json:"sport_type,omitempty" bson:"sport_type"`               // May take one of the following values: cycling, running, triathlon, other
	City            string `json:"city,omitempty" bson:"city"`                           // The club's city.
	State           string `json:"state,omitempty" bson:"state"`                         // The club's state or geographical region.
	Country         string `json:"country,omitempty" bson:"country"`                     // The club's country.
	Private         bool   `json:"private,omitempty" bson:"private"`                     // Whether the club is private.
	MemberCount     int    `json:"member_count,omitempty" bson:"member_count"`           // The club's member count.
	Featured        bool   `json:"featured,omitempty" bson:"featured"`                   // Whether the club is featured or not.
	Verified        bool   `json:"verified,omitempty" bson:"verified"`                   // Whether the club is verified or not.
	Url             string `json:"url,omitempty" bson:"url"`                             // The club's vanity URL.
	Membership      string `json:"membership,omitempty" bson:"membership"`               // The membership status of the logged-in athlete. May take one of the following values: member, pending
	Admin           bool   `json:"admin,omitempty" bson:"admin"`                         // Whether the currently logged-in athlete is an administrator of this club.
	Owner           bool   `json:"owner,omitempty" bson:"owner"`                         // Whether the currently logged-in athlete is the owner of this club.
	FollowingCount  int    `json:"following_count,omitempty" bson:"following_count"`     // The number of athletes in the club that the logged-in athlete follows.
}

// SubscriptionEvent 推送事件
type SubscriptionEvent struct {
	AspectType     string                 `json:"aspect_type" binding:"required"`
	EventTime      int64                  `json:"event_time" binding:"required"`
	ObjectID       int64                  `json:"object_id" binding:"required"`
	ObjectType     string                 `json:"object_type" binding:"required"`
	OwnerID        int64                  `json:"owner_id" binding:"required"`
	SubscriptionID int64                  `json:"subscription_id" binding:"required"`
	Updates        map[string]interface{} `json:"updates" binding:"required"`
}
